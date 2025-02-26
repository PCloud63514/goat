package environment

import "log"

var env = newViperEnvironment()

func init() {
	if env == nil {
		log.Fatal("Environment must not be nil")
	}
}

func New() Environment {
	return newViperEnvironment()
}

func GetProfiles() []string {
	return env.GetProfiles()
}

func ContainsProfile(expression string) bool {
	return env.ContainsProfile(expression)
}

func ContainsProperty(key string) bool {
	return env.ContainsProperty(key)
}

func GetPropertyString(key string, defaultValue string) string {
	return env.GetPropertyString(key, defaultValue)
}

func GetPropertyInt(key string, defaultValue int) int {
	return env.GetPropertyInt(key, defaultValue)
}

func GetPropertyBool(key string, defaultValue bool) bool {
	return env.GetPropertyBool(key, defaultValue)
}

func GetPropertySlice(key string, defaultValue []interface{}) []interface{} {
	return env.GetPropertySlice(key, defaultValue)
}

func GetRequiredPropertyString(key string) (string, error) {
	return env.GetRequiredPropertyString(key)
}

func GetRequiredPropertyInt(key string) (int, error) {
	return env.GetRequiredPropertyInt(key)
}

func GetRequiredPropertyBool(key string) (bool, error) {
	return env.GetRequiredPropertyBool(key)
}

func GetRequiredPropertySlice(key string) ([]interface{}, error) {
	return env.GetRequiredPropertySlice(key)
}
