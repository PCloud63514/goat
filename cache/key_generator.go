package cache

import "fmt"

type keyGenerator struct{}

func (*keyGenerator) Generate(params ...any) string {
	key := ""
	for _, param := range params {
		key += fmt.Sprintf("%v-", param)
	}
	return key[:len(key)-1]
}
