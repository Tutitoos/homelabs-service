package domain

const (
	DNSStatusNameOnline    string = "DNS – Conexión en línea"
	DNSStatusNameOffline   string = "DNS – Conexión perdida"
	DNSStatusNameRecovered string = "DNS – Conexión recuperada"
)

const (
	DNSStatusDescOnline    string = "🟢 En línea"
	DNSStatusDescOffline   string = "🔴 Caído"
	DNSStatusDescRecovered string = "🟡 Recuperado"
)

var DNSStatusNameMap = map[int]string{
	1: DNSStatusNameOnline,
	2: DNSStatusNameOffline,
	3: DNSStatusNameRecovered,
}

var DNSStatusDescMap = map[int]string{
	1: DNSStatusDescOnline,
	2: DNSStatusDescOffline,
	3: DNSStatusDescRecovered,
}

const (
	DNSCloudflarePrimary   string = "1.1.1.1"
	DNSCloudflareSecondary string = "1.0.0.1"
	DNSGooglePrimary       string = "8.8.8.8"
	DNSGoogleSecondary     string = "8.8.4.4"
)

const (
	DNSDescCloudflarePrimary   string = "Cloudflare Primary DNS"
	DNSDescCloudflareSecondary string = "Cloudflare Secondary DNS"
	DNSDescGooglePrimary       string = "Google Primary DNS"
	DNSDescGoogleSecondary     string = "Google Secondary DNS"
)

var DNSZoneMap = map[int]string{
	1: DNSCloudflarePrimary,
	2: DNSCloudflareSecondary,
	3: DNSGooglePrimary,
	4: DNSGoogleSecondary,
}

var DNSDescZoneMap = map[int]string{
	1: DNSDescCloudflarePrimary,
	2: DNSDescCloudflareSecondary,
	3: DNSDescGooglePrimary,
	4: DNSDescGoogleSecondary,
}

type DNSStruct struct {
}

func NewDNS() *DNSStruct {
	return &DNSStruct{}
}

func (d *DNSStruct) IsValidStatusId(statusId int) bool {
	_, exists := DNSStatusDescMap[statusId]
	return exists
}

func (d *DNSStruct) GetStatusName(statusId int) string {
	return DNSStatusNameMap[statusId]
}

func (d *DNSStruct) GetStatusDesc(statusId int) string {
	return DNSStatusDescMap[statusId]
}

func (d *DNSStruct) IsValidDNSId(dnsId int) bool {
	_, exists := DNSZoneMap[dnsId]
	return exists
}

func (d *DNSStruct) GetDNSName(dnsId int) string {
	return DNSZoneMap[dnsId]
}

func (d *DNSStruct) GetDNSDesc(dnsId int) string {
	return DNSDescZoneMap[dnsId]
}

var DNS = NewDNS()
