package datasources

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"homelabs-service/src/shared"

	"go.mongodb.org/mongo-driver/v2/bson"

	mongo "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoSortOptions struct {
	Field string
	Order int
}

type MongoCollections struct {
	SAI *mongo.Collection
	DNS *mongo.Collection
}

type IMongoDataSource struct {
	Client      *mongo.Client
	Collections MongoCollections
}

var (
	mongoClientInstance *mongo.Client
	mongoClientOnce     sync.Once
)

func MongoDataSource() *IMongoDataSource {
	return &IMongoDataSource{Client: nil, Collections: MongoCollections{}}
}

func (r *IMongoDataSource) Connect() (*IMongoDataSource, error) {
	ctx := context.Background()

	var err error
	mongoClientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(shared.Config.DatabaseMongoDbHost)
		clientOptions.SetAuth(options.Credential{
			Username: shared.Config.DatabaseMongoDbUser,
			Password: shared.Config.DatabaseMongoDbPassword,
		})

		mongoClientInstance, err = mongo.Connect(clientOptions)
		if err != nil {
			return
		}

		if errPing := mongoClientInstance.Ping(ctx, nil); errPing != nil {
			err = fmt.Errorf("error pinging MongoDB: %w", errPing)
			return
		}

	})

	if err != nil {
		return nil, err
	}

	r.Client = mongoClientInstance
	r.PreloadCollections()

	return r, nil
}

func (r *IMongoDataSource) PreloadCollections() {
	dbHomelabs := r.Client.Database("homelabs")

	r.Collections = MongoCollections{
		SAI: dbHomelabs.Collection("sai"),
		DNS: dbHomelabs.Collection("dns"),
	}
}

func (r *IMongoDataSource) GetCollection(name string) *mongo.Collection {
	switch name {
	case "sai":
		return r.Collections.SAI
	case "dns":
		return r.Collections.DNS
	default:
		return nil
	}
}

func (r *IMongoDataSource) Ping(ctx context.Context) (int64, error) {
	start := time.Now()
	if err := r.Client.Ping(ctx, nil); err != nil {
		return 0, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	return time.Since(start).Milliseconds(), nil
}

func (r *IMongoDataSource) Find(ctx context.Context, collectionName string, filter bson.M, limit int, page int, sortOptions MongoSortOptions) ([]bson.M, error) {
	shared.CapturePanic()

	// Defensive checks to avoid nil receiver panics
	if r == nil {
		return nil, fmt.Errorf("mongo data source is nil")
	}

	if r.Client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", collectionName)
	}

	optionsCollection := options.Find()
	optionsCollection.SetLimit(int64(limit))
	optionsCollection.SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2,
	})

	if page > 0 {
		optionsCollection.SetSkip(int64((page - 1) * limit))
	}

	if sortOptions != (MongoSortOptions{}) {
		optionsCollection.SetSort(bson.D{
			{Key: sortOptions.Field, Value: sortOptions.Order},
		})
	}

	if filter == nil {
		filter = bson.M{}
	}

	if limit == 0 {
		limit = 100
	}

	cursor, err := collection.Find(ctx, filter, optionsCollection)
	if err != nil {
		return nil, fmt.Errorf("error finding documents in %s: %w", collectionName, err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if errCursor := cursor.All(ctx, &results); errCursor != nil {
		return nil, fmt.Errorf("error decoding documents in %s: %w", collectionName, errCursor)
	}
	for i := range results {
		if id, ok := results[i]["_id"].(bson.ObjectID); ok {
			results[i]["_id"] = id.Hex()
		}
	}

	return results, nil
}

func (r *IMongoDataSource) FindOne(ctx context.Context, collectionName string, filter interface{}) (bson.M, error) {
	shared.CapturePanic()

	// Check if the client is nil
	if r == nil {
		return nil, fmt.Errorf("mongoDataSource is nil")
	}

	if r.Client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", collectionName)
	}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error finding one document in %s: %w", collectionName, err)
	}

	if id, ok := result["_id"].(bson.ObjectID); ok {
		result["_id"] = id.Hex()
	}

	return result, nil
}

func (r *IMongoDataSource) Create(ctx context.Context, collectionName string, document interface{}) (interface{}, error) {
	shared.CapturePanic()

	// Defensive checks to avoid nil receiver panics
	if r == nil {
		return nil, fmt.Errorf("mongo data source is nil")
	}

	if r.Client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return nil, fmt.Errorf("collection %s not found", collectionName)
	}

	var newDocument bson.M
	if v, ok := document.(bson.M); ok {
		newDocument = v
	} else {
		jsonMap, err := bson.Marshal(document)
		if err != nil {
			return nil, fmt.Errorf("error marshalling document: %w", err)
		}
		if err := bson.Unmarshal(jsonMap, &newDocument); err != nil {
			return nil, fmt.Errorf("error unmarshalling document: %w", err)
		}
	}

	delete(newDocument, "_id")

	result, err := collection.InsertOne(ctx, newDocument)
	if err != nil {
		return nil, fmt.Errorf("error creating document in %s: %w", collectionName, err)
	}

	if result.InsertedID == nil {
		return nil, fmt.Errorf("not created document in %s: %v", collectionName, document)
	}

	return result.InsertedID, nil
}

func (r *IMongoDataSource) UpdateOne(ctx context.Context, collectionName string, filter interface{}, document interface{}) error {
	shared.CapturePanic()

	// Defensive checks to avoid nil receiver panics
	if r == nil {
		return fmt.Errorf("mongo data source is nil")
	}

	if r.Client == nil {
		return fmt.Errorf("mongo client is nil")
	}

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	var newDocument bson.M
	if v, ok := document.(bson.M); ok {
		newDocument = v
	} else {
		jsonMap, err := bson.Marshal(document)
		if err != nil {
			return fmt.Errorf("error marshalling update document: %w", err)
		}
		if err := bson.Unmarshal(jsonMap, &newDocument); err != nil {
			return fmt.Errorf("error unmarshalling update document: %w", err)
		}
	}

	delete(newDocument, "_id")

	result, err := collection.UpdateOne(ctx, filter, newDocument)
	if err != nil {
		return fmt.Errorf("error updating document in %s: %w", collectionName, err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("not updated document in %s: %v", collectionName, filter)
	}

	return nil
}

func (r *IMongoDataSource) DeleteOne(ctx context.Context, collectionName string, filter interface{}) error {
	shared.CapturePanic()

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting document in %s: %w", collectionName, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not deleted document in %s: %v", collectionName, filter)
	}

	return nil
}

func (r *IMongoDataSource) DeleteMany(ctx context.Context, collectionName string, filter interface{}) error {
	shared.CapturePanic()

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting documents in %s: %w", collectionName, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not deleted document in %s: %v", collectionName, filter)
	}

	return nil
}

func (r *IMongoDataSource) Count(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	shared.CapturePanic()

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return 0, fmt.Errorf("collection %s not found", collectionName)
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("error counting documents in %s: %w", collectionName, err)
	}

	return count, nil
}

func (r *IMongoDataSource) CheckExists(ctx context.Context, collectionName string, filter bson.M) (bool, error) {
	shared.CapturePanic()

	collection := r.GetCollection(collectionName)
	if collection == nil {
		return false, fmt.Errorf("collection %s not found", collectionName)
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error counting documents in %s: %w", collectionName, err)
	}

	return count > 0, nil
}

func (r *IMongoDataSource) CreateIndex(ctx context.Context, collectionName string, keys bson.D, unique bool) error {
	shared.CapturePanic()

	collection := r.GetCollection(collectionName)
	if collection == nil {
		fmt.Printf("collection %s not found\n", collectionName)
		return nil
	}

	indexOptions := options.Index().SetUnique(unique).SetName(getIndexName(keys))

	if unique {
		partialFilter := bson.M{}

		for _, key := range keys {
			// Si es un campo anidado como "spotify.id"
			parts := strings.Split(key.Key, ".")
			if len(parts) > 1 {
				// Asegurar que el objeto raíz (ej. spotify) no sea null
				partialFilter[parts[0]] = bson.M{
					"$exists": true,
					"$type":   "string",
				}
			}

			// Validar que la propiedad final (ej. spotify.id) exista, sea string y no esté vacía
			partialFilter[key.Key] = bson.M{
				"$exists": true,
				"$type":   "string",
			}
		}

		indexOptions.SetPartialFilterExpression(partialFilter)
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions,
	}

	exists, _ := r.indexExists(ctx, collection, keys)
	if exists {
		fmt.Printf("Index already exists on collection %s with keys: %v\n", collectionName, keys)
		return nil
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			fmt.Printf("Duplicate keys found on collection %s\n", collectionName)
			return nil
		}

		fmt.Printf("Error creating index on collection %s: %v\n", collectionName, err)
		return nil
	}

	fmt.Printf("Index created successfully on collection %s with keys: %v\n", collectionName, keys)
	return nil
}

func (r *IMongoDataSource) indexExists(ctx context.Context, collection *mongo.Collection, keys bson.D) (bool, error) {
	shared.CapturePanic()

	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	var existing bson.M
	for cursor.Next(ctx) {
		if err := cursor.Decode(&existing); err != nil {
			continue
		}
		if keySpec, ok := existing["key"].(bson.M); ok {
			match := true
			for _, k := range keys {
				if keySpec[k.Key] != k.Value {
					match = false
					break
				}
			}
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func getIndexName(keys bson.D) string {
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s_%v", strings.ReplaceAll(k.Key, ".", "_"), k.Value))
	}
	return strings.Join(parts, "_") + "_custom"
}
