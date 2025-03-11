package environment

import (
	"fmt"
	"github.com/magiconair/properties"
	"regexp"
	"strconv"
	"sync"
)

const (
	FilePath  = "res"
	Extension = "properties"
)

var (
	sliceRegex *regexp.Regexp = regexp.MustCompile(`,\s*`)
)

type Environment struct {
	err     error
	mu      sync.RWMutex
	source  PropertySource
	sources map[string]PropertySource
}

type PropertySource struct {
	name     string
	resource map[string]string
}

func New(option ...string) *Environment {
	rootPropertySource := &PropertySource{
		name:     "",
		resource: map[string]string{},
	}
	sources := make(map[string]PropertySource)

	for _, opt := range option {
		if opt == "" {
			continue
		}
		resource := loadResource(opt)
		mergeMap(rootPropertySource.resource, resource)
		source := PropertySource{
			name:     opt,
			resource: resource,
		}
		sources[opt] = source
	}

	return &Environment{
		mu:      sync.RWMutex{},
		source:  *rootPropertySource,
		sources: sources,
	}
}

func (env *Environment) getProperty(key string) (string, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()
	if value, ok := env.source.resource[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("The [key=%s] property does not exist.", key)
}

func (env *Environment) ContainsProperty(key string) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()
	_, ok := env.source.resource[key]
	return ok
}

func (env *Environment) GetProperty(key string, value string) string {
	if v, err := env.getProperty(key); err == nil {
		return v
	}
	return value
}

func (env *Environment) GetPropertyInt(key string, value int) int {
	if v, err := env.getProperty(key); err == nil {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return value
}

func (env *Environment) GetPropertyBool(key string, value bool) bool {
	if v, err := env.getProperty(key); err == nil {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return value
}

func (env *Environment) GetPropertySlice(key string, value []string) []string {
	if v, err := env.getProperty(key); err == nil {
		return sliceRegex.Split(v, -1)
	}
	return value
}

func (env *Environment) GetRequiredProperty(key string) (string, error) {
	return env.getProperty(key)
}

func (env *Environment) GetRequiredPropertyInt(key string) (int, error) {
	v, err := env.getProperty(key)
	if err == nil {
		if i, err := strconv.Atoi(v); err == nil {
			return i, nil
		}
	}
	return 0, err
}

func (env *Environment) GetRequiredPropertyBool(key string) (bool, error) {
	v, err := env.getProperty(key)
	if err == nil {
		if b, err := strconv.ParseBool(v); err == nil {
			return b, nil
		}
	}
	return false, err
}

func (env *Environment) GetRequiredPropertySlice(key string) ([]string, error) {
	v, err := env.getProperty(key)
	if err == nil {
		return sliceRegex.Split(v, -1), nil
	}
	return []string{}, err
}

func (env *Environment) SetProperty(key string, value string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.source.resource[key] = value
}

func mergeMap(m1, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func loadResource(opt string) map[string]string {
	path := fmt.Sprintf("%s/%s.%s", FilePath, opt, Extension)
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}
	}
	return p.Map()
}
