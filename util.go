package siftjs

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// GetKind strips the pointer and interface from the kind
func GetKind(value interface{}) reflect.Kind {
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv.Kind()
}

// IsString check if the interface is a string
func IsString(value interface{}) bool {
	return GetKind(value) == reflect.String
}

// IsMap check if the interface is a map
func IsMap(value interface{}) bool {
	return GetKind(value) == reflect.Map
}

// IsArrayLike check if the interface is a slice or array
func IsArrayLike(value interface{}) bool {
	switch kind := GetKind(value); kind {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

// IsNumberLike check if interface is an int
func IsNumberLike(value interface{}) bool {
	switch kind := GetKind(value); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsIntLike check if interface is an int
func IsIntLike(value interface{}) bool {
	switch kind := GetKind(value); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

// Mapify turns an interface into a map
func Mapify(src interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if !IsMap(src) {
		return m, fmt.Errorf("invalid map")
	}

	if err := ToInterface(src, &m); err != nil {
		return m, err
	}

	return m, nil
}

// Arrayify turns an interface into a map
func Arrayify(src interface{}) ([]interface{}, error) {
	m := make([]interface{}, 0)

	// make array if not an array
	if !IsArrayLike(src) {
		a := make([]interface{}, 0)
		src = append(a, src)
	}

	if err := ToInterface(src, &m); err != nil {
		return m, err
	}

	return m, nil
}

// ToInterface converts one interface to another using json
func ToInterface(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &dest)
}

// FromJSON converts a json string to an interface
func FromJSON(str string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return nil
	}
	return result
}
