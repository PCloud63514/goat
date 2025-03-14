package environment

var (
	std = New()
)

func StandardEnvironment() *Environment {
	return std
}

func Reset(opts ...Option) {
	std = New(opts...)
}

func GetRequiredProperty(key string) (string, error) {
	return std.GetRequiredProperty(key)
}

func GetRequiredPropertyInt(key string) (int, error) {
	return std.GetRequiredPropertyInt(key)
}

func GetRequiredPropertyBool(key string) (bool, error) {
	return std.GetRequiredPropertyBool(key)
}

func GetRequiredPropertySlice(key string) ([]string, error) {
	return std.GetRequiredPropertySlice(key)
}

func GetProperty(key string, value string) string {
	return std.GetProperty(key, value)
}

func GetPropertyInt(key string, value int) int {
	return std.GetPropertyInt(key, value)
}

func GetPropertyBool(key string, value bool) bool {
	return std.GetPropertyBool(key, value)
}

func GetPropertySlice(key string, value []string) []string {
	return std.GetPropertySlice(key, value)
}

func Configuration(instance interface{}) (*interface{}, error) {
	return std.Configuration(instance)
}

func GetConfigurations() []interface{} {
	return std.configurations
}

func SetConfigurations(configurations []interface{}) {
	std.configurations = configurations
}

func SetProperty(key string, value string) {
	std.SetProperty(key, value)
}

func ContainsProperty(key string) bool {
	return std.ContainsProperty(key)
}

func GetKeys(prefix string) []string {
	return std.GetKeys(prefix)
}
