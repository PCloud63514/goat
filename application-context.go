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

type HandlerFunc func(ctx context.Context, env *Environment)

type GoatApplication struct {
	mu            sync.RWMutex
	startDateTime time.Time
	environment   Environment
	ctx           context.Context
	cancelFunc    context.CancelFunc
	chains        map[HandlerType][]HandlerFunc
}

func NewGoatApplication() *GoatApplication {
	return &GoatApplication{}
}

func (app *GoatApplication) Run() {
	app.startDateTime = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	app.ctx = ctx
	app.cancelFunc = cancel
	app.onInitialize()
	app.onStarting()
	app.onStarted()
	app.onPulling()
	app.onStop()
}

func (app *GoatApplication) onInitialize() {
	app.environment = *NewEnvironment()
}

func (app *GoatApplication) onStarting() {
	for _, execute := range app.getHandlers(HandlerType_Starting) {
		execute(app.ctx, &app.environment)
	}
}

func (app *GoatApplication) onStarted() {
	for _, execute := range app.getHandlers(HandlerType_Started) {
		execute(app.ctx, &app.environment)
	}
}

func (app *GoatApplication) onStop() {
	for _, execute := range app.getHandlers(HandlerType_Stop) {
		execute(app.ctx, &app.environment)
	}
}

func (app *GoatApplication) onPulling() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		app.cancelFunc()
	}()
	<-app.ctx.Done()
}

func (app *GoatApplication) getHandlers(t HandlerType) []HandlerFunc {
	if v, ok := app.chains[t]; ok {
		return v
	}
	return make([]HandlerFunc, 0)
}
