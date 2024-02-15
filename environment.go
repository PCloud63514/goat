package goat

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

const (
	PROFILE_DEFAULT    = "default"
	PROFILE_FLAG       = "profile"
	profile_SHORT_FLAG = "p"
)

/*
Profiles, PropertySources Management.
*/
type Environment interface {
	GetProfiles() []string
	ContainsProfile(expression string) bool
	GetPropertySources() map[string]interface{}
	PropertyResolver
}

type GoatEnvironment struct {
	Environment
	resolver PropertyResolver
	profiles []string
}

func NewGoatEnvironment(propertySource map[string]interface{}) *GoatEnvironment {
	defaultProfiles := []string{PROFILE_DEFAULT}
	readProfiles := readProfilesOfFlag()
	profiles := MergeSlicesUnique(defaultProfiles, readProfiles)

	return &GoatEnvironment{
		resolver: NewPropertySourcePropertyResolver(propertySource),
		profiles: profiles,
	}
}

func (env *GoatEnvironment) GetProfiles() []string {
	return env.profiles
}

func (env *GoatEnvironment) ContainsProfile(expression string) bool {
	for _, profile := range env.profiles {
		if expression == profile {
			return true
		}
	}
	return false
}

func (env *GoatEnvironment) GetPropertySource() map[string]interface{} {
	return env.resolver.GetPropertySource()
}

type PropertyResolver interface {
	GetPropertySource() map[string]interface{}
	ContainsProperty(key string) bool
	GetPropertyString(key string, defaultValue string) string
	GetPropertyInt(key string, defaultValue int) int
	GetPropertyBool(key string, defaultValue bool) bool
	GetRequiredPropertyString(key string) (string, error)
	GetRequiredPropertyInt(key string) (int, error)
	GetRequiredPropertyBool(key string) (bool, error)
}

type PropertySourcePropertyResolver struct {
	PropertyResolver
	propertySource map[string]interface{}
}

func NewPropertySourcePropertyResolver(propertySource map[string]interface{}) *PropertySourcePropertyResolver {
	_propertySource := propertySource
	if _propertySource == nil {
		_propertySource = make(map[string]interface{})
	}
	return &PropertySourcePropertyResolver{
		propertySource: _propertySource,
	}
}

func (p *PropertySourcePropertyResolver) GetPropertySource() map[string]interface{} {
	return p.propertySource
}

func (p *PropertySourcePropertyResolver) ContainsProperty(key string) bool {
	if _, err := p.getProperty(key); err != nil {
		return false
	}
	return true
}

func (p *PropertySourcePropertyResolver) GetPropertyString(key string, defaultValue string) string {
	value, err := p.getProperty(key)
	if err != nil {
		return defaultValue
	}
	return value.(string)
}

func (p *PropertySourcePropertyResolver) GetPropertyInt(key string, defaultValue int) int {
	value, err := p.getProperty(key)
	if err != nil {
		return defaultValue
	}
	i, err := strconv.Atoi(value.(string))
	if err != nil {
		return defaultValue
	}
	return i
}

func (p *PropertySourcePropertyResolver) GetPropertyBool(key string, defaultValue bool) bool {
	value, err := p.getProperty(key)
	if err != nil {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value.(string))
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func (p *PropertySourcePropertyResolver) GetRequiredPropertyString(key string) (string, error) {
	value, err := p.getProperty(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (p *PropertySourcePropertyResolver) GetRequiredPropertyInt(key string) (int, error) {
	value, err := p.getProperty(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(value.(string))
	if err != nil {
		return 0, fmt.Errorf("The [name=%s, value=%s] property is not Integer.", key, value)
	}
	return i, nil
}

func (p *PropertySourcePropertyResolver) GetRequiredPropertyBool(key string) (bool, error) {
	value, err := p.getProperty(key)
	if err != nil {
		return false, err
	}
	boolValue, err := strconv.ParseBool(value.(string))
	if err != nil {
		return false, fmt.Errorf("The [name=%s, value=%s] property is not Bool.", key, value)
	}
	return boolValue, nil
}

func (p *PropertySourcePropertyResolver) getProperty(key string) (interface{}, error) {
	pKey := p.formattedKey(key)
	if p.propertySource != nil {
		value, ok := p.propertySource[pKey]
		if !ok {
			return nil, fmt.Errorf("The [name=%s] property does not exist.", key)
		}
		return value, nil
	}
	return nil, fmt.Errorf("The PropertySource is null.")
}

func (p *PropertySourcePropertyResolver) formattedKey(key string) string {
	lowerKey := strings.ToLower(key)
	return lowerKey
}

func readProfilesOfFlag() []string {
	profiles := flag.String(PROFILE_FLAG, "", "Comma-separated list of profiles")
	p := flag.String(profile_SHORT_FLAG, "", "Comma-separated list of profiles (shorthand)")
	flag.Parse()

	if profiles != nil && *profiles != "" {
		return strings.Split(*profiles, ",")
	}

	if p != nil && *p != "" {
		return strings.Split(*p, ",")
	}
	return make([]string, 0)
}
