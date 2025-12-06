package shared

import (
	"fmt"
	"homelabs-service/src/domain/queries"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var DefaultBool = false
var DefaultString = ""
var DefaultInt = 0
var DefaultInt64 = int64(0)
var DefaultStringSlice = []string{}

type IParser struct {
}

type IParserStructure interface {
	SafeBool(b *bool) (bool, error)
	SafeString(s *string) (string, error)
	SafeInt(i *int) (int, error)
	SafeInt64(i *int64) (int64, error)
	SafeStringSlice(s *[]string) ([]string, error)
	FlattenMap(data bson.M, prefix string, out bson.M)
	ParseFormData(rawBody string) *queries.DNS
}

func Parser() IParserStructure {
	return &IParser{}
}

func (p *IParser) SafeBool(b *bool) (bool, error) {
	if b == nil {
		// log.Println("Error: SafeBool received a nil pointer", b)
		return DefaultBool, fmt.Errorf("nil pointer received in SafeBool value: %v", b)
	}

	if *b != true && *b != false {
		// log.Println("Warning: SafeBool received an invalid boolean value")
		return DefaultBool, fmt.Errorf("invalid boolean value received in SafeBool: %v", *b)
	}

	return *b, nil
}

func (p *IParser) SafeString(s *string) (string, error) {
	if s == nil {
		// log.Println("Error: SafeString received a nil pointer", s)
		return DefaultString, fmt.Errorf("nil pointer received in SafeString value: %v", s)
	}

	if *s == "" {
		// log.Println("Warning: SafeString received an empty string")
		return DefaultString, fmt.Errorf("empty string received in SafeString")
	}

	return *s, nil
}

func (p *IParser) SafeInt(i *int) (int, error) {
	if i == nil {
		// log.Println("Error: SafeInt received a nil pointer", i)
		return DefaultInt, fmt.Errorf("nil pointer received in SafeInt value: %v", i)
	}

	if *i < 0 {
		// log.Printf("Warning: SafeInt received a negative value: %d", *i)
		return DefaultInt, fmt.Errorf("negative value received in SafeInt: %d", *i)
	}

	return *i, nil
}

func (p *IParser) SafeInt64(i *int64) (int64, error) {
	if i == nil {
		// log.Println("Error: SafeInt64 received a nil pointer", i)
		return DefaultInt64, fmt.Errorf("nil pointer received in SafeInt64 value: %v", i)
	}

	if *i < 0 {
		// log.Printf("Warning: SafeInt64 received a negative value: %d", *i)
		return DefaultInt64, fmt.Errorf("negative value received in SafeInt64: %d", *i)
	}

	return *i, nil
}

func (p *IParser) SafeStringSlice(s *[]string) ([]string, error) {
	if s == nil {
		// log.Println("Error: SafeStringSlice received a nil pointer", s)
		return DefaultStringSlice, fmt.Errorf("nil pointer received in SafeStringSlice value: %v", s)
	}

	if len(*s) == 0 {
		// log.Println("Warning: SafeStringSlice received an empty slice")
		return DefaultStringSlice, fmt.Errorf("empty slice received in SafeStringSlice")
	}

	for idx, str := range *s {
		if str == "" {
			// log.Printf("Error: SafeStringSlice contains an empty string at index %d", idx)
			return DefaultStringSlice, fmt.Errorf("empty string found in SafeStringSlice at index %d", idx)
		}
	}

	return *s, nil
}

func (p *IParser) FlattenMap(data bson.M, prefix string, out bson.M) {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch child := v.(type) {
		case bson.M:
			p.FlattenMap(child, key, out)
		default:
			out[key] = v
		}
	}
}

func (p *IParser) ParseFormData(rawBody string) *queries.DNS {
	bodyData := new(queries.DNS)

	for _, pair := range strings.Split(rawBody, "&") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "dns_id":
			if dnsId, err := strconv.Atoi(val); err == nil {
				bodyData.DNSId = &dnsId
			}
		case "status_id":
			if statusId, err := strconv.Atoi(val); err == nil {
				bodyData.StatusId = &statusId
			}
		case "created_at":
			if createdAt, err := strconv.ParseInt(val, 10, 64); err == nil {
				bodyData.CreatedAt = &createdAt
			}
		}
	}

	return bodyData
}

var PARSER = Parser()
