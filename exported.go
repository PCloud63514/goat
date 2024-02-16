package goat

func Run(rType RunType) {
	app.Run(rType)
}

func AddHandlerFunc(hFunc HandlerFunc, t HandlerType) {
	app.AddHandlerFunc(hFunc, t)
}
