package application

import "fmt"

type applicationContext struct {
	environment
	components map[string]interface{}
}

type ApplicationContext interface {
	Environment
	RegisterComponent(name string, component interface{})
	GetComponent(name string) (interface{}, error)
	GetComponentAll() []interface{}
}

func New() *applicationContext {
	environment := newEnvironment()
	components := make(map[string]interface{})
	return &applicationContext{
		environment: *environment,
		components:  components,
	}
}

func (appCtx *applicationContext) RegisterComponent(name string, component interface{}) {
	appCtx.components[name] = component
}

func (appCtx *applicationContext) GetComponent(name string) (interface{}, error) {
	value, ok := appCtx.components[name]
	if !ok {
		return nil, fmt.Errorf("The %s component does not exist.")
	}
	return value, nil
}

func (appCtx *applicationContext) GetComponentAll() []interface{} {
	var instances []interface{}
	for _, instance := range appCtx.components {
		instances = append(instances, instance)
	}
	return instances
}
