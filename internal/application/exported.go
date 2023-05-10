package application

import (
	"fmt"
	"time"
)

var (
	appCtx ApplicationContext = New()
)

func GetProfile() string {
	return appCtx.GetProfile()
}

func GetGoVersion() string {
	return appCtx.GetGoVersion()
}

func GetPID() int {
	return appCtx.GetPID()
}

func GetStartUpDateTime() time.Time {
	return appCtx.GetStartUpDateTime()
}

func ContainsProperty(key string) bool {
	return appCtx.ContainsProperty(key)
}

func SetProperty(key string, value interface{}) {
	appCtx.SetProperty(key, value)
}

func GetProperty(key string) (interface{}, error) {
	return appCtx.GetProperty(key)
}

func GetPropertyString(key string) (string, error) {
	return appCtx.GetPropertyString(key)
}

func GetPropertyInt(key string) (int, error) {
	return appCtx.GetPropertyInt(key)
}

func GetPropertyBool(key string) (bool, error) {
	return appCtx.GetPropertyBool(key)
}

func GetOrDefaultProperty(key string, defaultValue interface{}) interface{} {
	return appCtx.GetOrDefaultProperty(key, defaultValue)
}

func GetOrDefaultPropertyString(key string, defaultValue string) string {
	return appCtx.GetOrDefaultPropertyString(key, defaultValue)
}

func GetOrDefaultPropertyInt(key string, defaultValue int) int {
	return appCtx.GetOrDefaultPropertyInt(key, defaultValue)
}

func GetOrDefaultPropertyBool(key string, defaultValue bool) bool {
	return appCtx.GetOrDefaultPropertyBool(key, defaultValue)
}

func RegisterComponent(name string, component interface{}) {
	appCtx.RegisterComponent(name, component)
}

func GetComponentAll() []interface{} {
	return appCtx.GetComponentAll()
}

func GetComponent[T any](name string) (T, error) {
	comp, err := appCtx.GetComponent(name)
	if nil != err {
		var t T
		return t, err
	}
	typedComp, ok := comp.(T)
	if !ok {
		var t T
		return t, fmt.Errorf("The Type of %s component is not %T", name, t)
	}
	return typedComp, nil
}

func OnError(err error) {
	if nil != err {
		panic(err)
	}
}
