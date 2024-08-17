package data

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"

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

// IDbCollection implementation
func (col *MongoCollection[T]) Upsert(document T) (models.IEntity, error) {
	if document.GetId() == "" {
		document.SetId(primitive.NewObjectID().Hex())
	}

	_, err := col.InsertOne(context.Background(), document)
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
