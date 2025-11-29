package repositories

import (
	"homelabs-service/src/infrastructure/datasources"
)

var (
	SAI ISAIMethodsRepository
)

func InitializeRepositories() {
	SAI = SAIRepository(datasources.DbMongo)
}

func init() {
	SAI = SAIRepository(nil)
}
