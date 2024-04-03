package goat

import "context"

func Run() {
	app.Run()
}

func AddHandlerFunc(hFunc HandlerFunc, t HandlerType) {
	app.AddHandlerFunc(hFunc, t)
}

func GetProfiles() []string {
	return app.environment.getProfiles()
}

func ContainsProfile(expression string) bool {
	return app.environment.containsProfile(expression)
}

func ContainsProperty(key string) bool {
	return app.environment.containsProperty(key)
}

func GetPropertyString(key string, defaultValue string) string {
	return app.environment.getPropertyString(key, defaultValue)
}

func GetPropertyInt(key string, defaultValue int) int {
	return app.environment.getPropertyInt(key, defaultValue)
}

func GetPropertyBool(key string, defaultValue bool) bool {
	return app.environment.getPropertyBool(key, defaultValue)
}

func GetRequiredPropertyString(key string) (string, error) {
	return app.environment.getRequiredPropertyString(key)
}

func GetRequiredPropertyInt(key string) (int, error) {
	return app.environment.getRequiredPropertyInt(key)
}

func GetRequiredPropertyBool(key string) (bool, error) {
	return app.environment.getRequiredPropertyBool(key)
}

func SetProperty(key string, value interface{}) {
	app.environment.setProperty(key, value)
}

func GetCtx() context.Context {
	return app.ctx
}

func GetCtxCancelFunc() context.CancelFunc {
	return app.cancelFunc
}

func OnError(err error) {
	if nil != err {
		panic(err)
	}
}
