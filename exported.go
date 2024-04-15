package goat

import "context"

func Run() {
	app.Run()
}

func AddHandlerFunc(hFunc HandlerFunc, t HandlerType) {
	app.AddHandlerFunc(hFunc, t)
}

func GetProfiles() []string {
	return app.environment.GetProfiles()
}

func ContainsProfile(expression string) bool {
	return app.environment.ContainsProfile(expression)
}

func ContainsProperty(key string) bool {
	return app.environment.ContainsProperty(key)
}

func GetPropertyString(key string, defaultValue string) string {
	return app.environment.GetPropertyString(key, defaultValue)
}

func GetPropertyInt(key string, defaultValue int) int {
	return app.environment.GetPropertyInt(key, defaultValue)
}

func GetPropertyBool(key string, defaultValue bool) bool {
	return app.environment.GetPropertyBool(key, defaultValue)
}

func GetRequiredPropertyString(key string) (string, error) {
	return app.environment.GetRequiredPropertyString(key)
}

func GetRequiredPropertyInt(key string) (int, error) {
	return app.environment.GetRequiredPropertyInt(key)
}

func GetRequiredPropertyBool(key string) (bool, error) {
	return app.environment.GetRequiredPropertyBool(key)
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
