package goat

import (
	"github.com/PCloud63514/goat/environment"
	"github.com/PCloud63514/goat/profile"
	"reflect"
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
	env := environment.New(prof.Get()...)

	// hook 기본 생성해야함.
	// opts로 전달받은 정보도 hook에 추가

	return &Goat{
		startUpDateTime: startUpDateTime,
		profile:         prof,
		environment:     env,
	}
}

func (g *Goat) Run() {

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
