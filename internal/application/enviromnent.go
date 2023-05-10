package application

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Environment interface {
	GetProfile() string
	GetGoVersion() string
	GetPID() int
	GetStartUpDateTime() time.Time
	ContainsProperty(key string) bool
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	GetPropertyString(key string) (string, error)
	GetPropertyInt(key string) (int, error)
	GetPropertyBool(key string) (bool, error)
	GetOrDefaultProperty(key string, defaultValue interface{}) interface{}
	GetOrDefaultPropertyString(key string, defaultValue string) string
	GetOrDefaultPropertyInt(key string, defaultValue int) int
	GetOrDefaultPropertyBool(key string, defaultValue bool) bool
}

type environment struct {
	profile         string
	goVersion       string
	pid             int
	startUpDateTime *time.Time
	mu              sync.Mutex
	propertySource  map[string]interface{}
}

func newEnvironment() *environment {
	startUpDateTime := time.Now()
	profile := readProfileOfFlag()
	pid := os.Getpid()
	goVersion := runtime.Version()
	propertySource := settingPropertySource(profile)

	return &environment{
		profile:         profile,
		goVersion:       goVersion,
		pid:             pid,
		startUpDateTime: &startUpDateTime,
		propertySource:  propertySource,
	}
}

func settingPropertySource(profile string) map[string]interface{} {
	viper.AddConfigPath(PROPERTY_PATH)
	viper.SetConfigType(PROPERTY_FILE_TYPE)
	viper.SetConfigName(PROFILE_DEFAULT)
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if PROFILE_DEFAULT != profile {
		viper.SetConfigName(profile)
		if err := viper.MergeInConfig(); err != nil {
			logrus.Warningf("The env file does not exists inbound a readable profile:[%v]", profile)
		}
	}
	return viper.AllSettings()
}

func (env *environment) GetProfile() string {
	return env.profile
}

func (env *environment) GetGoVersion() string {
	return env.goVersion
}

func (env *environment) GetPID() int {
	return env.pid
}

func (env *environment) GetStartUpDateTime() time.Time {
	return *env.startUpDateTime
}

func (env *environment) ContainsProperty(key string) bool {
	lowerKey := strings.ToLower(key)
	_, ok := env.propertySource[lowerKey]
	return ok
}

func (env *environment) SetProperty(key string, value interface{}) {
	lowerKey := strings.ToLower(key)
	env.propertySource[lowerKey] = value
}

func (env *environment) GetProperty(key string) (interface{}, error) {
	lowerKey := strings.ToLower(key)
	value, ok := env.propertySource[lowerKey]
	if !ok {
		return nil, fmt.Errorf("The %s property does not exist.", key)
	}
	return value, nil
}

func (env *environment) GetPropertyString(key string) (string, error) {
	value, err := env.GetProperty(key)
	if nil != err {
		return "", err
	}
	return value.(string), nil
}

func (env *environment) GetPropertyInt(key string) (int, error) {
	value, err := env.GetProperty(key)
	if nil != err {
		return 0, err
	}
	i, err := strconv.Atoi(value.(string))
	if nil != err {
		return 0, fmt.Errorf("The value [%s] of %s property is not Integer.", value, key)
	}
	return i, nil
}

func (env *environment) GetPropertyBool(key string) (bool, error) {
	value, err := env.GetProperty(key)
	if err != nil {
		return false, err
	}

	boolValue, err := strconv.ParseBool(value.(string))
	if err != nil {
		return false, fmt.Errorf("The value [%s] of %s property is not Bool.", value, key)
	}
	return boolValue, nil
}

func (env *environment) GetOrDefaultProperty(key string, defaultValue interface{}) interface{} {
	if value, err := env.GetProperty(key); nil == err {
		return value
	}
	return defaultValue
}

func (env *environment) GetOrDefaultPropertyString(key string, defaultValue string) string {
	if value, err := env.GetPropertyString(key); nil == err {
		return value
	}
	return defaultValue
}

func (env *environment) GetOrDefaultPropertyInt(key string, defaultValue int) int {
	if value, err := env.GetPropertyInt(key); nil == err {
		return value
	}
	return defaultValue
}

func (env *environment) GetOrDefaultPropertyBool(key string, defaultValue bool) bool {
	if value, err := env.GetPropertyBool(key); nil == err {
		return value
	}
	return defaultValue
}

func readProfileOfFlag() string {
	profileFlag := flag.String(PROFILE_OPTION, PROFILE_DEFAULT, "프로필")
	flag.Parse()
	return *profileFlag
}

