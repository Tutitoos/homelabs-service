package dtos

import (
	"homelabs-service/src/domain"
	"homelabs-service/src/domain/queries"
	"homelabs-service/src/shared"
)

type Backup struct {
	ZoneId    int     `json:"zoneId"`
	ZoneName  *string `json:"zoneName"`
	Message   *string `json:"message"`
	CreatedAt int64   `json:"createdAt"`
}

func NewBackup(backup queries.Backup) Backup {
	zoneId, _ := shared.PARSER.SafeInt(backup.ZoneId)
	zoneName := domain.BACKUP.GetZoneName(zoneId)
	message, _ := shared.PARSER.SafeString(backup.Message)
	createdAt, _ := shared.PARSER.SafeInt64(backup.CreatedAt)

	return Backup{
		ZoneId:    zoneId,
		ZoneName:  &zoneName,
		Message:   &message,
		CreatedAt: createdAt,
	}
}
