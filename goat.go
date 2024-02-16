package goat

var (
	app *GoatApplication
)

func init() {
	app = NewGoatApplication()
}

func Run() {
	app.Run()
}
