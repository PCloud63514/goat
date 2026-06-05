package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MergeSlicesUnique(a, b []string) []string {
	m := make(map[string]bool)
	var result []string

	for _, item := range a {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}

func MakeFile(content any, fileName string) {
	file, err := os.Create("./" + fileName)
	if nil != err {
		panic(err)
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, content); nil != err {
		if nil != err {
			panic(err)
		}
	}
}

func ToUint(v interface{}) uint {
	switch value := v.(type) {
	case string:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return uint(uintValue)
		}
	case uint:
		return value
	default:
		fmt.Println("i is of a different type")
	}
	return 0
}

func StrPtr(s string) *string {
	return &s
}

func ParseInt(s string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return fallback
	}
	return value
}

func ParseInt8(s string, fallback int8) int8 {
	value, err := strconv.ParseInt(strings.TrimSpace(s), 10, 8)
	if err != nil {
		return fallback
	}
	return int8(value)
}

func ParseInt16(s string, fallback int16) int16 {
	value, err := strconv.ParseInt(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return fallback
	}
	return int16(value)
}

func ParseInt32(s string, fallback int32) int32 {
	value, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return fallback
	}
	return int32(value)
}

func ParseInt64(s string, fallback int64) int64 {
	value, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return fallback
	}
	return value
}

func ParseUint(s string, fallback uint) uint {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, strconv.IntSize)
	if err != nil {
		return fallback
	}
	return uint(value)
}

func ParseUint8(s string, fallback uint8) uint8 {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, 8)
	if err != nil {
		return fallback
	}
	return uint8(value)
}

func ParseUint16(s string, fallback uint16) uint16 {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return fallback
	}
	return uint16(value)
}

func ParseUint32(s string, fallback uint32) uint32 {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	if err != nil {
		return fallback
	}
	return uint32(value)
}

func ParseUint64(s string, fallback uint64) uint64 {
	value, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return fallback
	}
	return value
}

func ParseFloat32(s string, fallback float32) float32 {
	value, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	if err != nil {
		return fallback
	}
	return float32(value)
}

func ParseFloat64(s string, fallback float64) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return fallback
	}
	return value
}

func ParseBool(s string, fallback bool) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return fallback
	}
	return value
}
