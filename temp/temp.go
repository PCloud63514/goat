package temp

import "sync"

type Environment interface {
	PropertyResolver
	GetProfiles() []string
	ContainsProfile(expression string) bool
	GetPropertySources() *PropertySources
}

type GoatEnvironment struct {
	Environment
	PropertyResolver
	sources  PropertySources
	profiles []string
}

func NewGoatEnvironment() *GoatEnvironment {
	// profiles 가져옴
	// sources 가져옴
	return &GoatEnvironment{
		PropertyResolver: PropertySourcePropertyResolver{},
	}
}

type PropertyResolver interface {
	GetPropertySources() *PropertySources
	ContainsProperty(key string) bool
	GetPropertyString(key string, defaultValue string) string
	GetPropertyInt(key string, defaultValue int) int
	GetPropertyBool(key string, defaultValue bool) bool
	GetRequiredPropertyString(key string) (string, error)
	GetRequiredPropertyInt(key string) (int, error)
	GetRequiredPropertyBool(key string) (bool, error)
}

type PropertySources interface {
	Get(name string) *PropertySource
	Contains(name string) bool
	AddFirst(source PropertySource)
	AddLast(source PropertySource)
	Remove(name string)
	Replace(name string, source PropertySource)
}

type PropertySourcePropertyResolver struct {
	PropertyResolver
	sources *PropertySources
}

type PropertySource struct {
	name     string
	resource map[string]interface{}
}

func (p *PropertySource) Copy() PropertySource {
	newResource := make(map[string]interface{})
	for key, value := range p.resource {
		newResource[key] = value
	}
	return PropertySource{
		name:     p.name,
		resource: newResource,
	}
}

type MutablePropertySources struct {
	PropertySources
	mu      sync.RWMutex
	sources []PropertySource
}

func (p *MutablePropertySources) Get(name string) *PropertySource {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.sources != nil {
		for _, source := range p.sources {
			if source.name == name {
				_source := source.Copy()
				return &_source
			}
		}
	}
	return nil
}

func (p *MutablePropertySources) Contains(name string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.sources != nil {
		for _, source := range p.sources {
			if source.name == name {
				return true
			}
		}
	}
	return false
}
