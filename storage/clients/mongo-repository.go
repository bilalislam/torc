package clients

import (
	"context"
	"errors"
	"github.com/bilalislam/torc/storage"
	"github.com/bilalislam/torc/storage/models"
	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoRepository interface {
	storage.IRepository
	GetByQuery(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	GetByQueryWithOption(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error)
	GetOneByQuery(ctx context.Context, filter interface{}, model models.IModel) error
	UpdateByQuery(ctx context.Context, filter interface{}, update interface{}) (int64, error)
	DeleteByQuery(ctx context.Context, filter interface{}) (int64, error)
	CountByQuery(ctx context.Context, filter interface{}) (int64, error)
}

type MongoRepository struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

type MongoClientOptions struct {
	ConnectionString string
	DbName           string
	Collection       string
}

func (mr *MongoRepository) Connect(ctx context.Context, dbContext interface{}) error {
	mongoClientOptions := dbContext.(MongoClientOptions)
	nrMon := nrmongo.NewCommandMonitor(nil)
	clientOptions := options.Client().ApplyURI(mongoClientOptions.ConnectionString).SetMonitor(nrMon)
	client, err := mongo.Connect(ctx, clientOptions)
	collection := client.Database(mongoClientOptions.DbName).Collection(mongoClientOptions.Collection)

	mr.Client = client
	mr.Collection = collection

	return err
}

func (mr *MongoRepository) GetById(ctx context.Context, id string, model models.IModel) error {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return mr.GetOneByQuery(ctx, bson.M{"_id": objectId}, model)
}

func (mr *MongoRepository) GetByQuery(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {

	if mr.Collection == nil || mr.Client == nil {
		return nil, errors.New("mongo client or connection nil. Please initialize first")
	}
	cur, sr := mr.Collection.Find(ctx, filter)

	if sr != nil {
		return nil, sr
	}
	return cur, nil
}

func (mr *MongoRepository) CountByQuery(ctx context.Context, filter interface{}) (int64, error) {
	if mr.Collection == nil || mr.Client == nil {
		return 0, errors.New("mongo client or connection nil. Please initialize first")
	}
	count, sr := mr.Collection.CountDocuments(ctx, filter)

	if sr != nil {
		return 0, sr
	}
	return count, nil
}

func (mr *MongoRepository) GetByQueryWithOption(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	if mr.Collection == nil || mr.Client == nil {
		return nil, errors.New("mongo client or connection nil. Please initialize first")
	}
	cur, sr := mr.Collection.Find(ctx, filter, options)

	if sr != nil {
		return nil, sr
	}
	return cur, nil
}

func (mr *MongoRepository) GetOneByQuery(ctx context.Context, filter interface{}, model models.IModel) error {

	if mr.Collection == nil || mr.Client == nil {
		return errors.New("mongo client or connection nil. Please initialize first")
	}
	sr := mr.Collection.FindOne(ctx, filter)

	if sr.Err() != nil {
		return sr.Err()
	}
	err := sr.Decode(model)
	return err
}

func (mr *MongoRepository) Save(ctx context.Context, model models.IModel) error {
	if mr.Collection == nil || mr.Client == nil {
		return errors.New("mongo client or connection nil. Please initialize first")
	}
	_, err := mr.Collection.InsertOne(ctx, model)
	return err
}

func (mr *MongoRepository) UpdateByQuery(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	if mr.Collection == nil || mr.Client == nil {
		return 0, errors.New("mongo client or connection nil. Please initialize first")
	}
	result, err := mr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

func (mr *MongoRepository) DeleteByQuery(ctx context.Context, filter interface{}) (int64, error) {
	if mr.Collection == nil || mr.Client == nil {
		return 0, errors.New("mongo client or connection nil. Please initialize first")
	}
	result, err := mr.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (mr *MongoRepository) Delete(ctx context.Context, id string) (int64, error) {
	if mr.Collection == nil || mr.Client == nil {
		return 0, errors.New("mongo client or connection nil. Please initialize first")
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	res, err := mr.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (mr *MongoRepository) Update(ctx context.Context, model models.IModel) (int64, error) {
	if mr.Collection == nil || mr.Client == nil {
		return 0, errors.New("mongo client or connection nil. Please initialize first")
	}
	objectId, err := primitive.ObjectIDFromHex(model.GetId())
	sr, err := mr.Collection.ReplaceOne(ctx, bson.M{"_id": objectId}, model)
	if err != nil {
		return 0, err
	}

	return sr.ModifiedCount, err
}
