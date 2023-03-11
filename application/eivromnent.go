package application

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

/*
+---------------+
|  Environment  |
+---------------+
*/

type configureEnvironment struct {
	propertySource    map[string]interface{}
	systemProperties  map[string]interface{}
	systemEnvironment map[string]interface{}
}

func newConfigureEnvironment(arguments *applicationArguments) *configureEnvironment {
	propertySource := makePropertySource(arguments)
	systemProperties := make(map[string]interface{})
	systemEnvironment := make(map[string]interface{})
	return &configureEnvironment{
		propertySource:    propertySource,
		systemProperties:  systemProperties,
		systemEnvironment: systemEnvironment,
	}
}
func (env *configureEnvironment) ContainsProperty(key string) bool {
	lowerKey := strings.ToLower(key)
	_, ok := env.propertySource[lowerKey]
	return ok
}
func (env *configureEnvironment) SetProperty(key string, value interface{}) {
	lowerKey := strings.ToLower(key)
	env.propertySource[lowerKey] = value
}
func (env *configureEnvironment) GetProperty(key string) interface{} {
	lowerKey := strings.ToLower(key)
	value := env.propertySource[lowerKey]
	return value
}
func (env *configureEnvironment) GetPropertyString(key string) string {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]
	if !ok {
		return ""
	}
	return value.(string)
}
func (env *configureEnvironment) GetPropertyInt(key string) int {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(value.(string))
	if nil != err {
		panic(err)
	}
	return i
}
func (env *configureEnvironment) GetOrDefaultProperty(key string, defaultValue interface{}) interface{} {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]

	if ok {
		return value
	}
	return defaultValue
}
func (env *configureEnvironment) GetOrDefaultPropertyInt(key string, defaultValue int) int {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]
	if ok {
		i, err := strconv.Atoi(value.(string))
		if nil != err {
			panic(err)
		}
		return i
	}
	return defaultValue
}
func (env *configureEnvironment) GetOrDefaultPropertyString(key string, defaultValue string) string {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]

	if ok {
		return value.(string)
	}
	return defaultValue
}
func makePropertySource(arguments *applicationArguments) map[string]interface{} {
	viper.AddConfigPath("./res")
	viper.SetConfigType("env")
	viper.SetConfigName("default")
	viper.AutomaticEnv()
	if FlagProfileValue != arguments.Profile {
		viper.SetConfigName(arguments.Profile)
		if err := viper.MergeInConfig(); err != nil {
			log.WithFields(log.Fields{
				"Profile": arguments.Profile,
			}).Warning("The env file does not exists ina readable profile")
		}
	}

	propertySource := viper.AllSettings()
	return propertySource
}
