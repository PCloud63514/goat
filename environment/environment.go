package environment

type Environment interface {
	GetProfiles() []string
	ContainsProfile(expression string) bool
	ContainsProperty(key string) bool
	GetPropertyString(key string, defaultValue string) string
	GetPropertyInt(key string, defaultValue int) int
	GetPropertyBool(key string, defaultValue bool) bool
	GetRequiredPropertyString(key string) (string, error)
	GetRequiredPropertyInt(key string) (int, error)
	GetRequiredPropertyBool(key string) (bool, error)
	SetProperty(key string, value interface{})
}
