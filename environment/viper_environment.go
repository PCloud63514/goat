package environment

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	propertyFilePath      = "res"
	propertyFileExtension = "env"
)

var (
	sliceRegex          = regexp.MustCompile(`,\s*`)
	emptyPropertySource = propertySource{
		name:     "",
		resource: map[string]interface{}{},
	}
	systemPropertySource = propertySource{
		name:     profileSystem,
		resource: map[string]interface{}{},
	}
)

type viperEnvironment struct {
	mu              sync.RWMutex
	readProfiles    []string
	defaultProfiles []string
	sources         []propertySource
}

type propertySource struct {
	name     string
	resource map[string]interface{}
}

func newViperEnvironment() Environment {
	env := &viperEnvironment{
		mu:              sync.RWMutex{},
		readProfiles:    readProfiles(),
		defaultProfiles: defaultProfiles,
		sources:         []propertySource{},
	}
	env.addLastPropertySource(emptyPropertySource)
	env.addLastPropertySource(systemPropertySource)
	for _, profile := range env.GetProfiles() {
		if profile != "" {
			resource := readPropertySource(profile)
			source := propertySource{
				name:     profile,
				resource: resource,
			}
			env.addLastPropertySource(source)
		}
	}
	return env
}

func (env *viperEnvironment) GetProfiles() []string {
	return mergeSlicesUnique(env.readProfiles, env.defaultProfiles)
}

func (env *viperEnvironment) ContainsProfile(expression string) bool {
	for _, profile := range env.GetProfiles() {
		if expression == profile {
			return true
		}
	}
	return false
}

func (env *viperEnvironment) ContainsProperty(key string) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()

	pKey := env.formattedKey(key)
	if env.sources != nil {
		for _, source := range env.sources {
			if _, ok := source.resource[pKey]; ok {
				return true
			}
		}
		return false
	}

	return false
}

func (env *viperEnvironment) GetPropertyString(key string, defaultValue string) string {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		return value.(string)
	}
	return defaultValue
}

func (env *viperEnvironment) GetPropertyInt(key string, defaultValue int) int {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		if i, err := strconv.Atoi(value.(string)); err == nil {
			return i
		}
	}
	return defaultValue
}

func (env *viperEnvironment) GetPropertyBool(key string, defaultValue bool) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		if b, err := strconv.ParseBool(value.(string)); err == nil {
			return b
		}
	}
	return defaultValue
}

func (env *viperEnvironment) GetPropertySlice(key string, defaultValue []interface{}) []interface{} {
	env.mu.RLock()
	defer env.mu.RUnlock()
	if value, err := env.getProperty(key); err == nil {
		return toSlice(value.(string))
	}
	return defaultValue
}

func (env *viperEnvironment) GetRequiredPropertyString(key string) (string, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	value, err := env.getProperty(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (env *viperEnvironment) GetRequiredPropertyInt(key string) (int, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	value, err := env.getProperty(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(value.(string))
	if err != nil {
		return 0, fmt.Errorf("The [name=%s, value=%s] property is not Integer.", key, value)
	}

	return i, nil
}

func (env *viperEnvironment) GetRequiredPropertyBool(key string) (bool, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	value, err := env.getProperty(key)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(value.(string))
	if err != nil {
		return false, fmt.Errorf("The [name=%s, value=%s] property is not Bool.", key, value)
	}
	return b, nil
}

func (env *viperEnvironment) GetRequiredPropertySlice(key string) ([]interface{}, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	value, err := env.getProperty(key)
	if err != nil {
		return nil, err
	}
	return toSlice(value.(string)), nil
}

func toSlice(value string) []interface{} {
	stringSlice := sliceRegex.Split(value, -1)
	result := make([]interface{}, len(stringSlice))
	for i, v := range stringSlice {
		result[i] = v
	}
	return result
}

func (env *viperEnvironment) SetProperty(key string, value interface{}) {
	env.mu.Lock()
	defer env.mu.Unlock()
	pKey := env.formattedKey(key)
	if env.sources == nil {
		env.sources = []propertySource{}
	}

	if len(env.sources) <= 0 {
		emptyPropertySource := propertySource{
			name:     "",
			resource: map[string]interface{}{},
		}
		env.addLastPropertySource(emptyPropertySource)
	}
	env.sources[0].resource[pKey] = value
}

func (env *viperEnvironment) addFirstPropertySource(source propertySource) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	env.sources = append([]propertySource{source}, env.sources...)
}

func (env *viperEnvironment) addLastPropertySource(source propertySource) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	env.sources = append(env.sources, source)
}

func (env *viperEnvironment) removePropertySource(name string) {
	env.mu.RLock()
	defer env.mu.RUnlock()
	_idx := -1

	for idx, s := range env.sources {
		if s.name == name {
			_idx = idx
		}
	}

	if _idx >= 0 || _idx < len(env.sources) {
		env.sources = append(env.sources[:_idx], env.sources[_idx+1:]...)
	}
}

func (env *viperEnvironment) replacePropertySource(name string, source propertySource) {
	env.mu.RLock()
	defer env.mu.RUnlock()
	_idx := -1

	for idx, s := range env.sources {
		if s.name == name {
			_idx = idx
		}
	}

	if _idx >= 0 || _idx < len(env.sources) {
		_sources := env.sources[:_idx]
		_sources = append(_sources, source)
		env.sources = append(_sources, env.sources[_idx+1:]...)
	}
}

func (env *viperEnvironment) getProperty(key string) (interface{}, error) {
	pKey := env.formattedKey(key)
	if env.sources != nil {
		for _, source := range env.sources {
			if value, ok := source.resource[pKey]; ok {
				return value, nil
			}
		}
		return nil, fmt.Errorf("The [name=%s] property does not exist.", key)
	}
	return nil, fmt.Errorf("The PropertySources is null.")
}

func (env *viperEnvironment) formattedKey(key string) string {
	lowerKey := strings.ToLower(key)
	return lowerKey
}

func loadProperties(name string) map[string]interface{} {
	return nil
}

func readPropertySource(profile string) map[string]interface{} {
	v := viper.New()
	v.AddConfigPath(propertyFilePath)
	v.SetConfigType(propertyFileExtension)
	v.SetConfigName(profile)
	v.AutomaticEnv()
	v.ReadInConfig()
	return v.AllSettings()
}

func mergeSlicesUnique(a, b []string) []string {
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
