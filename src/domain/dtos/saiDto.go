package dtos

import (
	"homelabs-service/src/domain"
	"homelabs-service/src/domain/queries"
	"homelabs-service/src/shared"
)

type SAI struct {
	ZoneId     int     `json:"zoneId"`
	ZoneName   *string `json:"zoneName"`
	StatusId   int     `json:"statusId"`
	StatusName *string `json:"statusName"`
	StatusDesc *string `json:"statusDesc"`
	CreatedAt  int64   `json:"createdAt"`
}

func NewSAI(sai queries.SAI) SAI {
	zoneId, _ := shared.PARSER.SafeInt(sai.ZoneId)
	zoneName := domain.SAI.GetZoneName(zoneId)
	statusId, _ := shared.PARSER.SafeInt(sai.StatusId)
	statusName := domain.SAI.GetStatusName(statusId)
	statusDesc := domain.SAI.GetStatusDesc(statusId)
	createdAt, _ := shared.PARSER.SafeInt64(sai.CreatedAt)

	return SAI{
		ZoneId:     zoneId,
		ZoneName:   &zoneName,
		StatusId:   statusId,
		StatusName: &statusName,
		StatusDesc: &statusDesc,
		CreatedAt:  createdAt,
	}
}
