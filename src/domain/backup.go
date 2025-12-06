package domain

const (
	BackupZoneHomelabs string = "Homelabs"
	BackupZone4Mans    string = "4Mans"
)

var BackupZoneMap = map[int]string{
	1: BackupZoneHomelabs,
	2: BackupZone4Mans,
}

type backupHelper struct{}

var BACKUP = backupHelper{}

func (backupHelper) IsValidZoneId(zoneId int) bool {
	_, exists := BackupZoneMap[zoneId]
	return exists
}

func (backupHelper) GetZoneName(zoneId int) string {
	return BackupZoneMap[zoneId]
}
