package goat

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"sync"
)

const (
	PROFILE_DEFAULT         = "default"
	PROFILE_SYSTEM          = "system"
	PROFILE_FLAG            = "profile"
	PROFILE_SHORT_FLAG      = "p"
	PROFILE_SEP             = ","
	PROPERTY_FILE_PATH      = "res"
	PROPERTY_FILE_EXTENSION = "env"
)

type Environment struct {
	mu       sync.RWMutex
	profiles []string
	sources  []PropertySource
}

type PropertySource struct {
	name     string
	resource map[string]interface{}
}

func NewEnvironment() *Environment {
	profiles := MergeSlicesUnique(readProfilesOfFlag(), []string{PROFILE_DEFAULT})
	env := &Environment{
		mu:       sync.RWMutex{},
		profiles: profiles,
		sources:  []PropertySource{},
	}
	emptyPropertySource := PropertySource{
		name:     "",
		resource: map[string]interface{}{},
	}
	env.AddLastPropertySource(emptyPropertySource)
	for _, profile := range profiles {
		if profile != "" {
			resource := readPropertySource(profile)
			source := PropertySource{
				name:     profile,
				resource: resource,
			}
			env.AddLastPropertySource(source)
		}
	}
	return env
}

func (env *Environment) GetProfiles() []string {
	return env.profiles
}

func (env *Environment) ContainsProfile(expression string) bool {
	for _, profile := range env.profiles {
		if expression == profile {
			return true
		}
	}
	return false
}

func (env *Environment) ContainsProperty(key string) bool {
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

func (env *Environment) GetPropertyString(key string, defaultValue string) string {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		return value.(string)
	}
	return defaultValue
}

func (env *Environment) GetPropertyInt(key string, defaultValue int) int {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		if i, err := strconv.Atoi(value.(string)); err == nil {
			return i
		}
	}
	return defaultValue
}

func (env *Environment) GetPropertyBool(key string, defaultValue bool) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()

	if value, err := env.getProperty(key); err == nil {
		if b, err := strconv.ParseBool(value.(string)); err == nil {
			return b
		}
	}
	return defaultValue
}

func (env *Environment) GetRequiredPropertyString(key string) (string, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	value, err := env.getProperty(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (env *Environment) GetRequiredPropertyInt(key string) (int, error) {
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

func (env *Environment) GetRequiredPropertyBool(key string) (bool, error) {
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

func (env *Environment) SetProperty(key string, value interface{}) {
	env.mu.Lock()
	defer env.mu.Unlock()
	pKey := env.formattedKey(key)
	if env.sources == nil {
		env.sources = []PropertySource{}
	}

	if len(env.sources) <= 0 {
		emptyPropertySource := PropertySource{
			name:     "",
			resource: map[string]interface{}{},
		}
		env.AddLastPropertySource(emptyPropertySource)
	}
	env.sources[0].resource[pKey] = value
}

func (env *Environment) AddFirstPropertySource(source PropertySource) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	env.sources = append([]PropertySource{source}, env.sources...)
}

func (env *Environment) AddLastPropertySource(source PropertySource) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	env.sources = append(env.sources, source)
}

func (env *Environment) RemovePropertySource(name string) {
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

func (env *Environment) ReplacePropertySource(name string, source PropertySource) {
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

func (env *Environment) getProperty(key string) (interface{}, error) {
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

func (env *Environment) formattedKey(key string) string {
	lowerKey := strings.ToLower(key)
	return lowerKey
}

func loadProperties(name string) map[string]interface{} {
	return nil
}

func readProfilesOfFlag() []string {
	profiles := flag.String(PROFILE_FLAG, "", "Comma-separated list of profiles")
	p := flag.String(PROFILE_SHORT_FLAG, "", "Comma-separated list of profiles (shorthand)")
	flag.Parse()

	if profiles != nil && *profiles != "" {
		_profiles := strings.Replace(*profiles, PROFILE_DEFAULT, "", 0)
		_profiles = strings.Replace(*profiles, PROFILE_SYSTEM, "", 0)
		return strings.Split(_profiles, ",")
	}

	if p != nil && *p != "" {
		_p := strings.Replace(*p, PROFILE_DEFAULT, "", 0)
		_p = strings.Replace(*p, PROFILE_SYSTEM, "", 0)
		return strings.Split(_p, PROFILE_SEP)
	}
	return make([]string, 0)
}

func readPropertySource(profile string) map[string]interface{} {
	v := viper.New()
	v.AddConfigPath(PROPERTY_FILE_PATH)
	v.SetConfigType(PROPERTY_FILE_EXTENSION)
	v.SetConfigName(profile)
	v.AutomaticEnv()
	v.ReadInConfig()
	return v.AllSettings()
}
