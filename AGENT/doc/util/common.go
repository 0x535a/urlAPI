package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func UnixTimeToString(ts int64) string {
	t := time.Unix(ts, 0)
	return fmt.Sprintf("%d年%d月%d日 %02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
	)
}

func DiffSlice(a, b []string) (added, removed []string) {
	setA := make(map[string]struct{})
	setB := make(map[string]struct{})

	for _, v := range a {
		setA[v] = struct{}{}
	}
	for _, v := range b {
		setB[v] = struct{}{}
	}

	for _, v := range b {
		if _, found := setA[v]; !found {
			added = append(added, v)
		}
	}
	for _, v := range a {
		if _, found := setB[v]; !found {
			removed = append(removed, v)
		}
	}
	return
}

func FindNthPreviousTime(hour, minute, second, n int, reletiveTime int64) time.Time {
	if hour < 1 || hour > 12 ||
		minute < 0 || minute > 59 ||
		second < 0 || second > 60 ||
		n < 1 {
		return TimeNow()
	}

	now := time.Unix(reletiveTime, 0)

	amHour := hour % 12
	pmHour := (hour % 12) + 12

	year, month, day := now.Date()
	todayAM := time.Date(year, month, day, amHour, minute, second, 0, time.Local)
	todayPM := time.Date(year, month, day, pmHour, minute, second, 0, time.Local)

	var lastOccurrence time.Time
	var daysToSubtract int

	if now.After(todayPM) {
		lastOccurrence = todayPM
		daysToSubtract = (n - 1) / 2
		if n%2 == 0 {
			lastOccurrence = todayAM
		}
	} else if now.After(todayAM) {
		lastOccurrence = todayAM
		daysToSubtract = (n - 1) / 2
		if n%2 == 0 {
			lastOccurrence = todayPM.AddDate(0, 0, -1)
		}
	} else {
		lastOccurrence = todayPM.AddDate(0, 0, -1)
		daysToSubtract = n / 2
		if n%2 == 1 {
			lastOccurrence = todayAM.AddDate(0, 0, -1)
		}
	}

	result := lastOccurrence.AddDate(0, 0, -daysToSubtract)

	return result
}

func IsComparable(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func TimeNow() time.Time {
	//return time.Date(2025, 9, 5, 20, 10, 0, 0, time.Local)
	return time.Now()
}

func NaiveLocalToNaiveUTC(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		time.UTC,
	)
}

func NaiveUTCToNaiveLocal(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		time.Local,
	)
}

func ReplaceNilWithZeroValue(ptr interface{}) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("argument must be pointer")
	}
	v = v.Elem()
	sanitizeValue(v)
}

func sanitizeValue(v reflect.Value) {
	if !v.IsValid() {
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		sanitizeValue(v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			sanitizeValue(field)
		}
	case reflect.Slice:
		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			}
		} else {
			for i := 0; i < v.Len(); i++ {
				sanitizeValue(v.Index(i))
			}
		}
	case reflect.Map:
		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.MakeMap(v.Type()))
			}
		} else {
			// Maps are tricky to sanitize in-place for keys/values if they are not pointers
			// For now, we only handle simple cases or skip
			for _, key := range v.MapKeys() {
				val := v.MapIndex(key)
				// Cannot sanitize map values directly if they are not pointers or addressable
				// This is a limitation of this simple recursive sanitizer
				if val.Kind() == reflect.Ptr || val.Kind() == reflect.Struct || val.Kind() == reflect.Slice {
					sanitizeValue(val)
				}
			}
		}
	case reflect.Float32, reflect.Float64:
		if v.CanSet() {
			val := v.Float()
			if val != val { // NaN check
				v.SetFloat(0)
			}
		}
	}
}

func GetPtr[T any](v T) *T {
	return &v
}

func Substr(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		runes = runes[:n]
	}
	return string(runes)
}

func RemoveLeadingZeros(s string) string {
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "0"
	}
	return s
}
