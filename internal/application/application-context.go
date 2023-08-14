package application

import (
	"fmt"
	"time"
)

var (
	startUpDateTime time.Time
	components      map[string]interface{}
)

func init() {
	startUpDateTime = time.Now()
	components = make(map[string]interface{})
}

func StartUpDateTime() time.Time {
	return startUpDateTime
}

func RegistComponent(name string, comp interface{}) {
	components[name] = comp
}
func GetComponent(name string) (interface{}, error) {
	value, ok := components[name]
	if !ok {
		return nil, fmt.Errorf("The %s component does not exist.")
	}
	return value, nil
}
func GetComponentAll() []interface{} {
	var comps []interface{}
	for _, comp := range components {
		comps = append(comps, comp)
	}
	return comps
}

func GetComponentByType[T any](name string) (T, error) {
	comp, err := GetComponent(name)
	if nil != err {
		var t T
		return t, err
	}
	v, ok := comp.(T)
	if !ok {
		var t T
		return t, fmt.Errorf("The Type of %s component is not %T", name, t)
	}
	return v, nil
}
func GetComponentAllByType[T any]() []T {
	comps := []T{}
	for _, comp := range components {
		if v, ok := comp.(T); ok {
			comps = append(comps, v)
		}
	}
	return comps
}
