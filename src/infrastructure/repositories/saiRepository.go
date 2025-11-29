package repositories

import (
	"context"
	"fmt"
	"strings"

	"homelabs-service/src/domain"
	"homelabs-service/src/domain/entities"
	"homelabs-service/src/infrastructure/datasources"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ISAIMethodsRepository interface {
	GetCollectionName() string
	GetDataSource() *datasources.IMongoDataSource
	GetMany(filter bson.M, limit int) ([]entities.SAI, error)
	Get(filter bson.M) (*entities.SAI, error)
	Create(item entities.SAI) (*entities.SAI, error)
	Update(filter bson.M, update bson.M) (*entities.SAI, error)
	Delete(filter bson.M) error
}

type ISAIRepository struct {
	DataSource     *datasources.IMongoDataSource
	CollectionName string
	Context        context.Context
}

func SAIRepository(dataSource *datasources.IMongoDataSource) ISAIMethodsRepository {
	return &ISAIRepository{
		DataSource:     dataSource,
		CollectionName: "sai",
		Context:        context.Background(),
	}
}

func (r *ISAIRepository) GetCollectionName() string {
	return r.CollectionName
}

func (r *ISAIRepository) GetDataSource() *datasources.IMongoDataSource {
	return r.DataSource
}

func (r *ISAIRepository) GetMany(filter bson.M, limit int) ([]entities.SAI, error) {
	if r.DataSource == nil {
		return nil, fmt.Errorf("dataSource is nil - connection not established")
	}

	items, err := r.DataSource.Find(r.Context, r.CollectionName, filter, limit, 0, datasources.MongoSortOptions{})
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return nil, nil
		}

		return nil, err
	}

	var results []entities.SAI
	for _, item := range items {
		bytes, err := bson.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("error marshalling bot config: %w", err)
		}

		var result entities.SAI
		if err := bson.Unmarshal(bytes, &result); err != nil {
			return nil, fmt.Errorf("error unmarshalling bot config: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *ISAIRepository) Get(filter bson.M) (*entities.SAI, error) {
	if r.DataSource == nil {
		return nil, fmt.Errorf("dataSource is nil - connection not established")
	}

	item, err := r.DataSource.FindOne(r.Context, r.CollectionName, filter)
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return nil, nil
		}

		return nil, err
	}

	bytes, err := bson.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error marshalling bot config: %w", err)
	}

	var result entities.SAI
	if err := bson.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling bot config: %w", err)
	}

	return &result, nil
}

func (r *ISAIRepository) Create(item entities.SAI) (*entities.SAI, error) {
	if r.DataSource == nil {
		return nil, fmt.Errorf("dataSource is nil - connection not established")
	}

	documentId, err := r.DataSource.Create(r.Context, r.CollectionName, item)
	if err != nil {
		return nil, fmt.Errorf("error creating bot config: %w", err)
	}

	result, err := r.Get(bson.M{"_id": documentId})
	if err != nil {
		return nil, fmt.Errorf("error retrieving created bot config: %w", err)
	}

	return result, nil
}

func (r *ISAIRepository) Update(filter bson.M, update bson.M) (*entities.SAI, error) {
	if r.DataSource == nil {
		return nil, fmt.Errorf("dataSource is nil - connection not established")
	}

	resultBefore, errResultBefore := r.Get(filter)
	if resultBefore == nil || errResultBefore != nil {
		return nil, fmt.Errorf("bot config not found or error retrieving it: %w", errResultBefore)
	}

	flatUpdate := bson.M{}
	domain.FlattenMap(resultBefore, update, "", flatUpdate)

	if err := r.DataSource.UpdateOne(r.Context, r.CollectionName, filter, bson.M{"$set": flatUpdate}); err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	if resultBefore.DocumentId != "" {
		objectId, err := bson.ObjectIDFromHex(resultBefore.DocumentId)
		if err != nil {
			return nil, fmt.Errorf("error converting DocumentId to ObjectID: %w", err)
		}

		filter = bson.M{"_id": objectId}
	}

	result, err := r.Get(filter)
	if err != nil {
		return nil, fmt.Errorf("error retrieving bot config: %w", err)
	}

	return result, nil
}

func (r *ISAIRepository) Delete(filter bson.M) error {
	if r.DataSource == nil {
		return fmt.Errorf("dataSource is nil - connection not established")
	}

	check, err := r.DataSource.CheckExists(r.Context, r.CollectionName, filter)
	if err != nil {
		return fmt.Errorf("error checking bot config: %w", err)
	}

	if !check {
		return fmt.Errorf("bot config with filter %v does not exist", filter)
	}

	err = r.DataSource.DeleteOne(r.Context, r.CollectionName, filter)
	if err != nil {
		return fmt.Errorf("error deleting bot config: %w", err)
	}

	return nil
}
