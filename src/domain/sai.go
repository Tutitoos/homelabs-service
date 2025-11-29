package domain

const (
	SAIStatusNameBattery  string = "SAI – Corte de Luz"
	SAIStatusNameOnline   string = "SAI – Electricidad Restaurada"
	SAIStatusNameCritical string = "SAI – Batería Baja"
	SAIStatusNameShutdown string = "SAI – Forzando Apagado"
	SAIStatusNameEvent    string = "SAI – Evento"
)

const (
	SAIStatusDescBattery  string = "🔴 Funcionando con batería"
	SAIStatusDescOnline   string = "🟢 En línea"
	SAIStatusDescCritical string = "🟠 Batería en nivel crítico"
	SAIStatusDescShutdown string = "⚫ Apagado inminente"
	SAIStatusDescEvent    string = "ℹ️ Evento"
)

var SaiStatusNameMap = map[int]string{
	1: SAIStatusNameBattery,
	2: SAIStatusNameOnline,
	3: SAIStatusNameCritical,
	4: SAIStatusNameShutdown,
	5: SAIStatusNameEvent,
}

var SaiStatusDescMap = map[int]string{
	1: SAIStatusDescBattery,
	2: SAIStatusDescOnline,
	3: SAIStatusDescCritical,
	4: SAIStatusDescShutdown,
	5: SAIStatusDescEvent,
}

const (
	SAIZoneHomelabs string = "Homelabs"
	SAIZone4Mans    string = "4Mans"
)

var SaiZoneMap = map[int]string{
	1: SAIZoneHomelabs,
	2: SAIZone4Mans,
}

type SAIStruct struct {
}

func NewSAI() *SAIStruct {
	return &SAIStruct{}
}

func (s *SAIStruct) IsValidStatusId(statusId int) bool {
	_, exists := SaiStatusDescMap[statusId]
	return exists
}

func (s *SAIStruct) GetStatusName(statusId int) string {
	return SaiStatusNameMap[statusId]
}

func (s *SAIStruct) GetStatusDesc(statusId int) string {
	return SaiStatusDescMap[statusId]
}

func (s *SAIStruct) IsValidZoneId(zoneId int) bool {
	_, exists := SaiZoneMap[zoneId]
	return exists
}

func (s *SAIStruct) GetZoneName(zoneId int) string {
	return SaiZoneMap[zoneId]
}

var SAI = NewSAI()
