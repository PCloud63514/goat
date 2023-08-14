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
)

const (
	profileDefault    = "default"
	profileSimpleFlag = "p"
	profileFlag       = "profile"
	profileExtension  = ".profile"

	pidExtension          = ".pid"
	propertyFileExtension = "env"

	makeFilePath         = "./"
	propertyFileReadPath = "res"
)

var (
	profile        = profileDefault
	propertySource = make(map[string]interface{})
)

func init() {
	profile = readProfileOfFlag()
	propertySource = readPropertySource(profile)

	go makeFile(PID(), pidExtension)
	go makeFile(profile, profileExtension)
}

func readPropertySource(profile string) map[string]interface{} {
	viper.AddConfigPath(propertyFileReadPath)
	viper.SetConfigType(propertyFileExtension)
	viper.SetConfigName(profileDefault)
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if profileDefault != profile {
		viper.SetConfigName(profile)
		if err := viper.MergeInConfig(); err != nil {
			logrus.Warningf("The env file does not exists inbound a readable profile:[%v]", profile)
		}
	}
	return viper.AllSettings()
}

func readProfileOfFlag() string {
	profileFlag := flag.String(profileFlag, profileDefault, "프로필")
	flag.Parse()
	return *profileFlag
}

func makeFile(content any, fileExtension string) {
	file, err := os.Create(makeFilePath + fileExtension)
	if nil != err {
		panic(err)
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, content); nil != err {
		if nil != err {
			panic(err)
		}
	}
}

func Profile() string {
	return profile
}

func GoVersion() string {
	return runtime.Version()
}

func PID() int {
	return os.Getpid()
}

func ContainsProperty(key string) bool {
	lowerKey := strings.ToLower(key)
	_, ok := propertySource[lowerKey]
	return ok
}

func SetProperty(key string, value interface{}) {
	lowerKey := strings.ToLower(key)
	propertySource[lowerKey] = value
}

func GetProperty(key string) (interface{}, error) {
	lowerKey := strings.ToLower(key)
	value, ok := propertySource[lowerKey]
	if !ok {
		return nil, fmt.Errorf("The %s property does not exist.", key)
	}
	return value, nil
}

func GetPropertyString(key string) (string, error) {
	value, err := GetProperty(key)
	if nil != err {
		return "", err
	}
	return value.(string), nil
}

func GetPropertyInt(key string) (int, error) {
	value, err := GetProperty(key)
	if nil != err {
		return 0, err
	}
	i, err := strconv.Atoi(value.(string))
	if nil != err {
		return 0, fmt.Errorf("The value [%s] of %s property is not Integer.", value, key)
	}
	return i, nil
}

func GetPropertyBool(key string) (bool, error) {
	value, err := GetProperty(key)
	if err != nil {
		return false, err
	}

	boolValue, err := strconv.ParseBool(value.(string))
	if err != nil {
		return false, fmt.Errorf("The value [%s] of %s property is not Bool.", value, key)
	}
	return boolValue, nil
}

func GetOrDefaultProperty(key string, defaultValue interface{}) interface{} {
	if value, err := GetProperty(key); nil == err {
		return value
	}
	return defaultValue
}

func GetOrDefaultPropertyString(key string, defaultValue string) string {
	if value, err := GetPropertyString(key); nil == err {
		return value
	}
	return defaultValue
}

func GetOrDefaultPropertyInt(key string, defaultValue int) int {
	if value, err := GetPropertyInt(key); nil == err {
		return value
	}
	return defaultValue
}

func GetOrDefaultPropertyBool(key string, defaultValue bool) bool {
	if value, err := GetPropertyBool(key); nil == err {
		return value
	}
	return defaultValue
}
