package data_models

type IDbCollection[TEntity any] interface {
	Upsert(collection TEntity, document TEntity) (TEntity, error)
	FirstOrDefault(collection TEntity, filter TEntity) TEntity
	Where(collection interface{}, filter interface{}) (TEntity, error)
}

type IDbProvider interface {
	Connect()
	Disconnect()
	GetClient() any
}
