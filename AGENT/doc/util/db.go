package util

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"strings"
)

func NewUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}

func NewToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

func IsEmpty(src interface{}) bool {
	if src == nil {
		return true
	}
	v := reflect.ValueOf(src)

	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		// 所有字段都是零值才算空
		for i := 0; i < v.NumField(); i++ {
			if !IsEmpty(v.Field(i).Interface()) {
				return false
			}
		}
		return true
	}
	return false
}

func StringListToDB(list []string) string {
	if len(list) == 0 {
		return ""
	}
	result := ""
	for i, item := range list {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%s", item)
	}
	return result
}

func DBToStringList(dbList string) []string {
	if dbList == "" {
		return nil
	}
	var result []string
	items := strings.Split(dbList, ",")
	for _, item := range items {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func RemoveFromDBList(dbList string, item string) string {
	if dbList == "" {
		return ""
	}
	items := strings.Split(dbList, ",")
	var result []string
	for _, i := range items {
		if i != item {
			result = append(result, i)
		}
	}
	return StringListToDB(result)
}

func ExistInDBList(dbList string, item string) bool {
	if dbList == "" {
		return false
	}
	items := strings.Split(dbList, ",")
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
