package environment

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/magiconair/properties"
)

const (
	defaultResourcePath = "res"
	defaultExtension    = "properties"
	defaultProfile      = "default"
	propertiesTagKey    = "properties"
	defaultTagKey       = "default"
	legacyDefaultPrefix = "default="
)

var (
	sliceRegex    *regexp.Regexp = regexp.MustCompile(`,\s*`)
	defaultOption Option         = Option{
		ResPath:        defaultResourcePath,
		Profiles:       []string{defaultProfile},
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
	applyConfigurationDefaults(rootPropertySource.resource, opt.Configurations)
	prop := properties.LoadMap(rootPropertySource.resource)
	for _, instance := range opt.Configurations {
		if err := prop.Decode(instance); err != nil {
			panic(err)
		}
	}

	return &Environment{
		mu:             sync.RWMutex{},
		prop:           prop,
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
	applyConfigurationDefaults(env.source.resource, []interface{}{instance})
	env.prop = properties.LoadMap(env.source.resource)
	if err := env.prop.Decode(instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

func (env *Environment) GetConfigurations() []interface{} {
	return env.configurations
}

func (env *Environment) SetConfigurations(configurations []interface{}) {
	applyConfigurationDefaults(env.source.resource, configurations)
	env.prop = properties.LoadMap(env.source.resource)
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
	path := fmt.Sprintf("%s/%s.%s", filePath, opt, defaultExtension)
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}
	}
	return p.Map()
}

func applyConfigurationDefaults(resource map[string]string, configurations []interface{}) {
	for _, configuration := range configurations {
		applyStructDefaults(resource, reflect.TypeOf(configuration))
	}
}

func applyStructDefaults(resource map[string]string, t reflect.Type) {
	if t == nil {
		return
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldType := field.Type
		resolvedFieldType := fieldType
		for resolvedFieldType.Kind() == reflect.Ptr {
			resolvedFieldType = resolvedFieldType.Elem()
		}

		if resolvedFieldType.Kind() == reflect.Struct {
			applyStructDefaults(resource, fieldType)
		}

		key, defaultValue := parsePropertyTags(field)
		if key == "" || defaultValue == nil {
			continue
		}
		if _, exists := resource[key]; exists {
			continue
		}

		resource[key] = *defaultValue
	}
}

func parsePropertyTags(field reflect.StructField) (string, *string) {
	tag := field.Tag.Get(propertiesTagKey)
	if tag == "" {
		return "", nil
	}

	parts := strings.Split(tag, ",")
	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", nil
	}

	if defaultTag, ok := field.Tag.Lookup(defaultTagKey); ok {
		return key, &defaultTag
	}

	return key, parseLegacyDefaultTag(parts[1:])
}

// parseLegacyDefaultTag는 기존 `properties:"key,default=value"` 형식과의 하위 호환을 위해 유지합니다.
// Deprecated: 새 코드는 `properties:"key" default:"value"` 형식을 사용해야 합니다.
func parseLegacyDefaultTag(parts []string) *string {
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, legacyDefaultPrefix) {
			value := strings.TrimPrefix(part, legacyDefaultPrefix)
			return &value
		}
	}

	return nil
}
