package goat

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type RunType string
type HandlerType string

const (
	RunType_Web          RunType     = "WEB"
	RunType_Standard     RunType     = "STANDARD"
	HandlerType_Starting HandlerType = "STARTING"
	HandlerType_Started  HandlerType = "STARTED"
	HandlerType_Stop     HandlerType = "STOP"
)

var (
	app *Goat = New()
)

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
	environment   *Environment
	chains        map[HandlerType][]HandlerFunc
	runType       RunType
	ctx           context.Context
	cancelFunc    context.CancelFunc
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
	app.environment = NewEnvironment()
	app.updateContext()
}

func (app *Goat) onStarting() {
	for _, execute := range app.getHandlers(HandlerType_Starting) {
		execute(app.ctx, app.environment)
	}
	app.makeFile(PID(), ".pid")
	app.makeFile(strings.Join(app.environment.readProfiles, ","), ".profile")
	app.applicationStartMsg()
}

func (app *Goat) onStarted() {
	for _, execute := range app.getHandlers(HandlerType_Started) {
		execute(app.ctx, app.environment)
	}
	// batch start
	StartScheduler()
	// if runType == web -> gin Listen
}

func (app *Goat) onStop() {
	StopScheduler()
	for _, execute := range app.getHandlers(HandlerType_Stop) {
		execute(app.ctx, app.environment)
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

func (app *Goat) applicationStartMsg() {
	endTime := time.Now()
	elapsedTime := endTime.Sub(app.startDateTime)

	logrus.WithFields(logrus.Fields{
		"StartupDateTime":  app.startDateTime.Format("2006-01-02 15:04:05"),
		"Profile":          strings.Join(app.environment.GetProfiles(), ","),
		"PID":              PID(),
		"GoVersion":        GoVersion(),
		"completedSeconds": fmt.Sprintf("%dm %ds", int(elapsedTime.Minutes()), int(elapsedTime.Seconds())%60),
	}).Info("Application Start!")
}

func GoVersion() string {
	return runtime.Version()
}

func PID() int {
	return os.Getpid()
}

func (goat *Goat) updateContext() {
	ctx, cancel := context.WithCancel(context.Background())
	goat.ctx = ctx
	goat.cancelFunc = cancel
	goat.runType = RunType(goat.environment.GetPropertyString("RUN_TYPE", string(RunType_Standard)))
}

func (goat *Goat) makeFile(content any, fileName string) {
	file, err := os.Create("./" + fileName)
	if nil != err {
		panic(err)
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, content); nil != err {
		if nil != err {
			panic(err)
		}
	}
}
