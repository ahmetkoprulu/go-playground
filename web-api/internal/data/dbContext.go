package data

import (
	"context"
	"errors"
	"math"
	"os"
	"time"

	data_models "github.com/ahmetkoprulu/go-playground/web-api/internal/data/abstract"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbContext *MongoDbContext
)

type MongoDbContext struct {
	*mongo.Database
}

type MongoCollection[T models.IEntity] struct {
	*mongo.Collection
}

func InitializeMongoDb() error {
	dbContext = &MongoDbContext{}
	err := dbContext.Connect()

	return err
}

func Context() (*MongoDbContext, error) {
	if dbContext.Database == nil {
		return nil, errors.New("MongoDbContext is not initialized")
	}

	return dbContext, nil
}

// Collections
func (ctx *MongoDbContext) Users() *MongoCollection[*models.User] {
	return &MongoCollection[*models.User]{
		ctx.Database.Collection("users"),
	}
}

func (ctx *MongoDbContext) Notifications() *MongoCollection[*models.Notification] {
	return &MongoCollection[*models.Notification]{
		ctx.Database.Collection("notifications"),
	}
}

// IDbCollection implementation
func (col *MongoCollection[T]) Upsert(document T) (models.IEntity, error) {
	if document.GetId() == "" {
		document.SetId(primitive.NewObjectID().Hex())
	}

	_, err := col.ReplaceOne(context.Background(), bson.M{"_id": document.GetId()}, document, options.Replace().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (col *MongoCollection[T]) FirstOrDefault(filter any) (T, error) {
	var result T
	err := col.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (col *MongoCollection[T]) Where(filter any) ([]T, error) {
	var result []T
	cursor, err := col.Find(context.Background(), filter)
	cursor.All(context.Background(), &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (col *MongoCollection[T]) Paginate(filter any, page int, take int) (data_models.PagingModel[T], error) {
	if filter == nil {
		filter = bson.M{}
	}

	var result data_models.PagingModel[T]
	skip := (page - 1) * take

	pipeline := mongo.Pipeline{{{Key: "$match", Value: filter}}, {{Key: "$facet", Value: bson.D{{Key: "metadata", Value: bson.A{bson.D{{Key: "$count", Value: "total"}}}}, {Key: "data", Value: bson.A{bson.D{{Key: "$skip", Value: skip}}, bson.D{{Key: "$limit", Value: take}}}}}}}}
	cursor, err := col.Aggregate(context.Background(), pipeline)
	if err != nil {
		return result, err
	}

	defer cursor.Close(context.Background())
	var aggregationResult data_models.AggregateResult[T]
	if err := cursor.All(context.Background(), &aggregationResult); err != nil {
		return result, err
	}

	if len(aggregationResult) > 0 && len(aggregationResult[0].Metadata) > 0 {
		result.TotalCount = aggregationResult[0].Metadata[0].Total
	}

	result.CurrentPage = page
	result.Take = take
	// ceil
	result.TotalPage = int(math.Ceil(float64(result.TotalCount) / float64(take)))
	result.Data = aggregationResult[0].Data

	return result, nil
}

func (col *MongoCollection[T]) Delete(id string) error {
	_, err := col.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// IDbProvider implementation
func (db *MongoDbContext) Connect() error {
	var connecionString = os.Getenv("CONNECT_STRING")
	if connecionString == "" {
		return errors.New("connection string is not provided")
	}

	clientOptions := options.Client().ApplyURI(connecionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	var dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		return errors.New("database name is not provided")
	}

	db.Database = client.Database(dbName)
	return nil
}

func (db *MongoDbContext) Disconnect() {
	ctx := context.Background()
	db.Database.Client().Disconnect(ctx)
	db.Database = nil
}

func (ctx *MongoDbContext) GetClient() any {
	return ctx.Database.Client()
}
