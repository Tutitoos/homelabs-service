package entities

import (
	"time"

	"homelabs-service/src/domain"
	"homelabs-service/src/domain/queries"
)

type DNS struct {
	DocumentId string `bson:"_id"`
	DNSId      int    `bson:"dns_id"`
	StatusId   int    `bson:"status_id"`
	CreatedAt  int64  `bson:"created_at"`
}

func CreateDNS(data queries.DNS) (*DNS, []string) {
	errors := []string{}

	if data.DNSId == nil {
		errors = append(errors, "dns_id is required")
	} else if !domain.DNS.IsValidDNSId(*data.DNSId) {
		errors = append(errors, "dns_id is invalid")
	}

	if data.StatusId == nil {
		errors = append(errors, "status_id is required")
	} else if !domain.DNS.IsValidStatusId(*data.StatusId) {
		errors = append(errors, "status_id is invalid")
	}

	// Si created_at no viene, generarlo automáticamente
	var createdAt int64
	if data.CreatedAt == nil {
		createdAt = time.Now().UnixMilli()
	} else {
		createdAt = *data.CreatedAt
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &DNS{
		DocumentId: "", // NOTE: DO NOT TOUCH THIS, LEAVE IT EMPTY. THE DATABASE WILL HANDLE ASSIGNING DOCUMENTID.
		DNSId:      *data.DNSId,
		StatusId:   *data.StatusId,
		CreatedAt:  createdAt,
	}, errors
}
