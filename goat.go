package goat

import (
	"context"
	"github.com/PCloud63514/goat/environment"
	"github.com/PCloud63514/goat/profile"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type Goat struct {
	err             error
	startUpDateTime time.Time
	profile         *profile.Profile
	environment     *environment.Environment
	hooks           map[HookType][]HookFunc
}

type Option interface{}

func New(opts ...Option) *Goat {
	startUpDateTime := time.Now()
	prof := profile.New()
	env := environment.New(environment.Option {
		Profiles: profile.Get(),
	})

	return &Goat{
		startUpDateTime: startUpDateTime,
		profile:         prof,
		environment:     env,
		hooks:           make(map[HookType][]HookFunc),
	}
}

func (g *Goat) Run() {
	if exitCode := g.run(g.Wait); exitCode != 0 {
		os.Exit(exitCode)
	}
}

func (g *Goat) Wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}

func (g *Goat) Start(ctx context.Context) (err error) {
	if g.err != nil {
		return g.err
	}
	return nil
}

func (g *Goat) Stop(ctx context.Context) (err error) {
	return nil
}

func (g *Goat) run(done func() <-chan os.Signal) (exitCode int) {
	startCtx, startCancel := context.WithCancel(context.Background())
	defer startCancel()

	if err := g.Start(startCtx); err != nil {
		return 1
	}

	<-done()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	defer stopCancel()

	if err := g.Stop(stopCtx); err != nil {
		return 1
	}

	return 0
}

func Provide(constructors ...interface{}) Option {
	for _, constructor := range constructors {
		fnType := reflect.TypeOf(constructor)
		if fnType.Kind() != reflect.Func {
			panic("constructor must be a function")
		}
	}

	// 그래프를 이 단계에서 만들어야하나? 아니다 지금은 그냥 확인만 하면서 추가만함.
	return nil
}

func Configuration(configurations ...interface{}) Option {
	return nil
}
