package repositories

import (
	"homelabs-service/src/infrastructure/datasources"
)

var (
	SAI ISAIMethodsRepository
	DNS IDNSMethodsRepository
)

func InitializeRepositories() {
	SAI = SAIRepository(datasources.DbMongo)
	DNS = DNSRepository(datasources.DbMongo)
}

func init() {
	SAI = SAIRepository(nil)
	DNS = DNSRepository(nil)
}
