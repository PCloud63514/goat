package application

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	initializeFunc []func()
	startFunc      []func()
	destroyFunc    []func(ctx context.Context)
)

type InitializeCommand interface {
	OnInitialize()
}
type DestroyCommand interface {
	OnDestroy()
}

func Run() {
	onInitialize()
	onStart()
	applicationStartMsg()
	fulling()
}

func onInitialize() {
	logrus.Trace("[LifeCycle] Initialize")
	initializer := GetComponentAllByType[InitializeCommand]()
	for _, command := range initializer {
		command.OnInitialize()
	}
	for _, fun := range initializeFunc {
		fun()
	}
	logrus.Trace("[LifeCycle] Initialize Complete!")

}

func onStart() {
	logrus.Trace("[LifeCycle] Start")
	for _, fun := range startFunc {
		fun()
	}
	logrus.Trace("[LifeCycle] Start Complete!")
}

func onDestroy(ctx context.Context) {
	logrus.Trace("[LifeCycle] Destroy")
	for _, fun := range destroyFunc {
		fun(ctx)
	}
	destroyer := GetComponentAllByType[DestroyCommand]()
	for _, command := range destroyer {
		command.OnDestroy()
	}
	logrus.Trace("[LifeCycle] Destroy Complete!")

}

func applicationStartMsg() {
	endTime := time.Now()
	elapsedTime := endTime.Sub(StartUpDateTime())
	logrus.WithFields(logrus.Fields{
		"StartupDateTime":  StartUpDateTime().Format("2006-01-02 15:04:05"),
		"Profile":          Profile(),
		"GoVersion":        GoVersion(),
		"PID":              PID(),
		"completedSeconds": fmt.Sprintf("%dm %ds", int(elapsedTime.Minutes()), int(elapsedTime.Seconds())%60),
	}).Info("Application Start!")
}

func fulling() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logrus.Info("shutdown...")
	onDestroy(ctx)
	logrus.Info("Shutdown complete")
	os.Exit(0)
}

func AddStartEventCallback(f func()) {
	startFunc = append(startFunc, f)
}

func AddDestroyEventCallback(f func(ctx context.Context)) {
	destroyFunc = append(destroyFunc, f)
}
