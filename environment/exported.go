package environment

func New() Environment {
	return newViperEnvironment()
}
