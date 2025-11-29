package datasources

import (
	"fmt"
	"sync"

	"homelabs-service/src/shared"
)

type IDataSources struct {
	MongoDataSource *IMongoDataSource
	MongoReady      chan bool
	mu              sync.Mutex
}

var (
	DbMongo *IMongoDataSource
)

func DataSources() *IDataSources {
	return &IDataSources{
		MongoReady: make(chan bool),
	}
}

func (r *IDataSources) Connect() {
	shared.CapturePanic()

	var wg sync.WaitGroup
	wg.Add(1)

	go r.MongoConnect()

	go func() {
		defer wg.Done()
		<-r.MongoReady
	}()

	wg.Wait()
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.MongoDataSource == nil {
		fmt.Println("❌ error connecting to MongoDB: MongoDataSource is nil")

		return
	}

	fmt.Println("✅ all repositories are connected")
}

func (r *IDataSources) MongoConnect() {
	dataSource, err := MongoDataSource().Connect()
	if err != nil {
		fmt.Printf("❌ error connecting to MongoDB: %v\n", err.Error())
		return
	}

	r.mu.Lock()
	r.MongoDataSource = dataSource
	DbMongo = dataSource
	r.mu.Unlock()

	r.MongoReady <- true

	fmt.Println("✅ connected to MongoDB")
}
