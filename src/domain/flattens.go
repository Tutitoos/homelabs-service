package domain

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func FlattenMap[T any](original T, data bson.M, prefix string, out bson.M) {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch child := v.(type) {
		case bson.M:
			if matchKey(original, key) {
				FlattenMap(original, child, key, out)
			} else {
				out[key] = child
			}
		default:
			out[key] = v
		}
	}
}

func matchKey[T any](original T, key string) bool {
	val := reflect.ValueOf(original).Elem()
	for _, part := range strings.Split(key, ".") {
		field := val.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, part)
		})
		if !field.IsValid() {
			return false
		}
		val = field
	}
	return !val.IsZero()
}
