package goat

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HandlerType string

const (
	HandlerType_Starting = "STARTING"
	HandlerType_Started  = "STARTED"
	HandlerType_Stop     = "STOP"
)

var (
	app *Goat = New()
)

func Run() {
	app.Run()
}

func New() *Goat {
	return &Goat{
		mu:     sync.RWMutex{},
		chains: make(map[HandlerType][]HandlerFunc),
	}
}

type HandlerFunc func(ctx context.Context, env *Environment)

type Goat struct {
	mu            sync.RWMutex
	startDateTime time.Time
	environment   Environment
	ctx           context.Context
	cancelFunc    context.CancelFunc
	chains        map[HandlerType][]HandlerFunc
}

func (app *Goat) Run() {
	app.onInitialize()
	app.onStarting()
	app.onStarted()
	app.onPulling()
	app.onStop()
}

func (app *Goat) onInitialize() {
	app.startDateTime = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	app.ctx = ctx
	app.cancelFunc = cancel
	app.environment = *NewEnvironment()
}

func (app *Goat) onStarting() {
	for _, execute := range app.getHandlers(HandlerType_Starting) {
		execute(app.ctx, &app.environment)
	}
}

func (app *Goat) onStarted() {
	for _, execute := range app.getHandlers(HandlerType_Started) {
		execute(app.ctx, &app.environment)
	}
}

func (app *Goat) onStop() {
	for _, execute := range app.getHandlers(HandlerType_Stop) {
		execute(app.ctx, &app.environment)
	}
}

func (app *Goat) onPulling() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		app.cancelFunc()
	}()
	<-app.ctx.Done()
}

func (app *Goat) getHandlers(t HandlerType) []HandlerFunc {
	if v, ok := app.chains[t]; ok {
		return v
	}
	return make([]HandlerFunc, 0)
}

func (app *Goat) AddHandlerFunc(hFunc HandlerFunc, t HandlerType) {
	if app.chains[t] == nil {
		app.chains[t] = []HandlerFunc{}
	}
	app.chains[t] = append(app.chains[t], hFunc)
}
