package dtos

import (
	"homelabs-service/src/domain"
	"homelabs-service/src/domain/entities"
	"homelabs-service/src/shared"
)

type DNS struct {
	DocumentId string  `json:"documentId"`
	DNSId      int     `json:"dnsId"`
	DNSName    *string `json:"dnsName"`
	StatusId   int     `json:"statusId"`
	StatusName *string `json:"statusName"`
	StatusDesc *string `json:"statusDesc"`
	CreatedAt  int64   `json:"createdAt"`
}

func NewDNS(dns entities.DNS) DNS {
	documentId, _ := shared.PARSER.SafeString(&dns.DocumentId)
	dnsId, _ := shared.PARSER.SafeInt(&dns.DNSId)
	dnsName := domain.DNS.GetDNSName(dnsId)
	statusId, _ := shared.PARSER.SafeInt(&dns.StatusId)
	statusName := domain.DNS.GetStatusName(statusId)
	statusDesc := domain.DNS.GetStatusDesc(statusId)
	createdAt, _ := shared.PARSER.SafeInt64(&dns.CreatedAt)

	return DNS{
		DocumentId: documentId,
		DNSId:      dnsId,
		DNSName:    &dnsName,
		StatusId:   statusId,
		StatusName: &statusName,
		StatusDesc: &statusDesc,
		CreatedAt:  createdAt,
	}
}
