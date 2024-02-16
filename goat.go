package goat

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type RunType string
type HandlerType string

const (
	RunType_Web          RunType     = "WEB"
	RunType_Default      RunType     = "DEFAULT"
	HandlerType_Starting HandlerType = "STARTING"
	HandlerType_Started  HandlerType = "STARTED"
	HandlerType_Stop     HandlerType = "STOP"
)

var (
	app *Goat = New()
)

func New() *Goat {
	return &Goat{
		mu:      sync.RWMutex{},
		chains:  make(map[HandlerType][]HandlerFunc),
		runType: RunType_Default,
	}
}

type HandlerFunc func(ctx context.Context, env *Environment)

type Goat struct {
	mu            sync.RWMutex
	startDateTime time.Time
	runType       RunType
	environment   Environment
	ctx           context.Context
	cancelFunc    context.CancelFunc
	chains        map[HandlerType][]HandlerFunc
}

func (app *Goat) Run(rType RunType) {
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
	// batch start
	// if runType == web -> gin Listen
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
