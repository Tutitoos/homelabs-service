package domain

const (
	StatusNameBattery  string = "SAI – Corte de Luz"
	StatusNameOnline   string = "SAI – Electricidad Restaurada"
	StatusNameCritical string = "SAI – Batería Baja"
	StatusNameShutdown string = "SAI – Forzando Apagado"
	StatusNameEvent    string = "SAI – Evento"
)

const (
	StatusDescBattery  string = "🔴 Funcionando con batería"
	StatusDescOnline   string = "🟢 En línea"
	StatusDescCritical string = "🟠 Batería en nivel crítico"
	StatusDescShutdown string = "⚫ Apagado inminente del servidor"
	StatusDescEvent    string = "ℹ️ Evento"
)

var SaiStatusNameMap = map[int]string{
	1: StatusNameBattery,
	2: StatusNameOnline,
	3: StatusNameCritical,
	4: StatusNameShutdown,
	5: StatusNameEvent,
}

var SaiStatusDescMap = map[int]string{
	1: StatusDescBattery,
	2: StatusDescOnline,
	3: StatusDescCritical,
	4: StatusDescShutdown,
	5: StatusDescEvent,
}

func IsValidStatusId(statusId int) bool {
	_, exists := SaiStatusDescMap[statusId]
	return exists
}

func GetStatusName(statusId int) string {
	return SaiStatusNameMap[statusId]
}

func GetStatusDesc(statusId int) string {
	return SaiStatusDescMap[statusId]
}

const (
	ZoneHomelabs string = "Homelabs"
	Zone4Mans    string = "4Mans"
)

var SaiZoneMap = map[int]string{
	1: ZoneHomelabs,
	2: Zone4Mans,
}

func IsValidZoneId(zoneId int) bool {
	_, exists := SaiZoneMap[zoneId]
	return exists
}

func GetZoneName(zoneId int) string {
	return SaiZoneMap[zoneId]
}
