package queries

type Backup struct {
	ZoneId    *int    `json:"zone_id"`
	Message   *string `json:"message"`
	CreatedAt *int64  `json:"created_at"`
}
