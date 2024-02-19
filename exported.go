package goat

func Run() {
	app.Run()
}

func AddHandlerFunc(hFunc HandlerFunc, t HandlerType) {
	app.AddHandlerFunc(hFunc, t)
}
