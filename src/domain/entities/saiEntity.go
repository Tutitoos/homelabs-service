package entities

import (
	"homelabs-service/src/domain"
	"homelabs-service/src/domain/queries"
)

type SAI struct {
	DocumentId string `bson:"_id"`
	ZoneId     int    `bson:"zone_id"`
	StatusId   int    `bson:"status_id"`
	CreatedAt  int64  `bson:"created_at"`
}

func CreateSAI(data queries.SAI) (*SAI, []string) {
	errors := []string{}

	if data.ZoneId == nil {
		errors = append(errors, "zone_id is required")
	} else if !domain.IsValidZoneId(*data.ZoneId) {
		errors = append(errors, "zone_id is invalid")
	}

	if data.StatusId == nil {
		errors = append(errors, "status_id is required")
	} else if !domain.IsValidStatusId(*data.StatusId) {
		errors = append(errors, "status_id is invalid")
	}

	if data.CreatedAt == nil {
		errors = append(errors, "created_at is required")
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &SAI{
		DocumentId: "", // NOTE: DO NOT TOUCH THIS, LEAVE IT EMPTY. THE DATABASE WILL HANDLE ASSIGNING DOCUMENTID.
		ZoneId:     *data.ZoneId,
		StatusId:   *data.StatusId,
		CreatedAt:  *data.CreatedAt,
	}, errors
}
