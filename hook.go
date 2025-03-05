package goat

type HookType int

type HookFunc func()

type HookOnStartFunc func()
type HookOnStopFunc func()
type HookOnShutdownFunc func()

const (
	hookType_Start HookType = iota
	hookType_Stop
	hookType_Shutdown
)

type Hook struct {
	OnStart    func() error
	OnStop     func() error
	OnShutdown func() error
}
