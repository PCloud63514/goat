package goat

import (
	"context"
	"fmt"
	"github.com/PCloud63514/goat/internal/utils"
	"github.com/PCloud63514/goat/logger"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type HandlerType string

const (
	HandlerType_Initialize HandlerType = "INITIALIZE"
	HandlerType_Starting   HandlerType = "STARTING"
	HandlerType_Started    HandlerType = "STARTED"
	HandlerType_Stop       HandlerType = "STOP"
	HandlerType_Destroy    HandlerType = "DESTROY"
)

var (
	app *Goat = New()
)

func New() *Goat {
	ctx, cancel := context.WithCancel(context.Background())
	return &Goat{
		mu:          sync.RWMutex{},
		chains:      make(map[HandlerType][]HandlerFunc),
		environment: newEnvironment(),
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}

type HandlerFunc func(ctx context.Context, env *Environment)

type Goat struct {
	mu               sync.RWMutex
	startRunDateTime time.Time
	environment      *Environment
	chains           map[HandlerType][]HandlerFunc
	ctx              context.Context
	cancelFunc       context.CancelFunc
}

func (app *Goat) Run() {
	app.onInitialize()
	app.onStarting()
	app.onStarted()
	app.onPulling()
	app.onStop()
}

func (app *Goat) onSystemInitialize() {
	app.startRunDateTime = time.Now()
	logger.Infof("Application Initialize\t Profile: [%v]\t PID: %v\tGoVersion: %v\t StartDateTime: %v\t",
		strings.Join(app.environment.GetProfiles(), ","),
		os.Getpid(),
		runtime.Version(),
		app.startRunDateTime.Format("2006-01-02 15:04:05"),
	)
	utils.MakeFile(os.Getpid(), ".pid")
	utils.MakeFile(strings.Join(app.environment.GetProfiles(), ","), ".profile")
}

func (app *Goat) onInitialize() {
	app.onSystemInitialize()
	app.executeHandlers(HandlerType_Initialize)
}

func (app *Goat) onStarting() {
	app.executeHandlers(HandlerType_Starting)
	app.applicationStartMsg()
}

func (app *Goat) onStarted() {
	app.executeHandlers(HandlerType_Started)
}

func (app *Goat) onPulling() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		app.cancelFunc()
	}()
	<-app.ctx.Done()
	logger.Info("shutdown...")
}

func (app *Goat) onStop() {
	app.executeHandlers(HandlerType_Stop)
}

func (app *Goat) onDestroy() {
	app.executeHandlers(HandlerType_Destroy)
	logger.Info("Shutdown complete")
}

func (app *Goat) applicationStartMsg() {
	endTime := time.Now()
	elapsedTime := endTime.Sub(app.startRunDateTime)
	logger.Infof("Application Start\n\t\t\t┠ Profile: [%v]\n\t\t\t┠ PID: %v\n\t\t\t┠ GoVersion: %v\n\t\t\t┠ StartDateTime: %v\n\t\t\t┖ CompletedSeconds: %v",
		strings.Join(app.environment.GetProfiles(), ","),
		os.Getpid(),
		runtime.Version(),
		app.startRunDateTime.Format("2006-01-02 15:04:05"),
		fmt.Sprintf("%dm %ds", int(elapsedTime.Minutes()), int(elapsedTime.Seconds())%60),
	)
	logger.Error(fmt.Errorf("ErrorTest"))
}

func (app *Goat) getHandlers(t HandlerType) []HandlerFunc {
	app.mu.RLock()
	defer app.mu.RUnlock()

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

func (app *Goat) executeHandlers(t HandlerType) {
	for _, execute := range app.getHandlers(t) {
		execute(app.ctx, app.environment)
	}
}
