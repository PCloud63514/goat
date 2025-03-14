package cache

import "fmt"

type KeyFunc func(params ...any) string

func defaultKeyFunc(params ...any) string {
	key := ""
	for _, param := range params {
		key += fmt.Sprintf("%v-", param)
	}
	return key[:len(key)-1] // 마지막 `-` 제거
}
