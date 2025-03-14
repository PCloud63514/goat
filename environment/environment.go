package environment

import (
	"fmt"
	"github.com/magiconair/properties"
	"regexp"
	"strconv"
	"sync"
)

const (
	defaultResourcePath = "res"
	Extension           = "properties"
)

var (
	sliceRegex    *regexp.Regexp = regexp.MustCompile(`,\s*`)
	defaultOption Option         = Option{
		ResPath:        defaultResourcePath,
		Profiles:       []string{},
		Configurations: make([]interface{}, 0),
	}
)

type Environment struct {
	err            error
	mu             sync.RWMutex
	source         PropertySource
	prop           *properties.Properties
	sources        map[string]PropertySource
	configurations []interface{}
	resourcePath   string
}

type Option struct {
	ResPath        string
	Profiles       []string
	Configurations []interface{}
}

func (opt *Option) apply(option Option) {
	if option.ResPath != "" {
		opt.ResPath = option.ResPath
	}
	if option.Profiles != nil && len(option.Profiles) > 0 {
		opt.Profiles = option.Profiles
	}
	if option.Configurations != nil && len(option.Configurations) > 0 {
		opt.Configurations = option.Configurations
	}
}

type PropertySource struct {
	name     string
	resource map[string]string
}

func New(opts ...Option) *Environment {
	opt := &Option{}
	opt.apply(defaultOption)
	for _, o := range opts {
		opt.apply(o)
	}
	return new(opt)
}

func new(opt *Option) *Environment {
	rootPropertySource := &PropertySource{
		name:     "",
		resource: map[string]string{},
	}
	sources := make(map[string]PropertySource)

	for _, profile := range opt.Profiles {
		if profile == "" {
			continue
		}
		resource := loadResource(opt.ResPath, profile)
		mergeMap(rootPropertySource.resource, resource)
		source := PropertySource{
			name:     profile,
			resource: resource,
		}
		sources[profile] = source
	}
	prop := properties.LoadMap(rootPropertySource.resource)
	for _, instance := range opt.Configurations {
		if err := prop.Decode(instance); err != nil {
			panic(err)
		}
	}
	return &Environment{
		mu:             sync.RWMutex{},
		prop:           properties.LoadMap(rootPropertySource.resource),
		source:         *rootPropertySource,
		sources:        sources,
		configurations: opt.Configurations,
	}
}

func (env *Environment) GetKeys(prefix string) []string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	keys := make([]string, 0)

	for key, _ := range env.source.resource {
		if prefix != "" && !regexp.MustCompile(prefix).MatchString(key) {
			continue
		}
		keys = append(keys, key)
	}

	return keys
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
	env.prop = properties.LoadMap(env.source.resource)
}

func (env *Environment) Configuration(instance interface{}) (*interface{}, error) {
	if err := env.prop.Decode(instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (env *Environment) GetConfigurations() []interface{} {
	return env.configurations
}

func (env *Environment) SetConfigurations(configurations []interface{}) {
	for _, instance := range configurations {
		if err := env.prop.Decode(instance); err != nil {
			panic(err)
		}
	}
	env.configurations = configurations
}

func mergeMap(m1, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func loadResource(filePath, opt string) map[string]string {
	path := fmt.Sprintf("%s/%s.%s", filePath, opt, Extension)
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}
	}
	return p.Map()
}
