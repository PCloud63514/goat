package utils

import (
	"fmt"
	"os"
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
