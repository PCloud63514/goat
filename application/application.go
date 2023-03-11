package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func Start() *ApplicationContext {
	startTime := time.Now()
	logger := newLogger()
	logger.Info("Application Startup")
	logger.Info("Application Context Loading")
	ctx := newConfigurableContext(logger)
	args := newApplicationArguments(startTime)
	env := newConfigureEnvironment(args)
	appCtx := newApplicationContext(ctx, args, env)
	router := gin.Default()

	port := appCtx.GetOrDefaultPropertyString("port", ":8080")

	if err := router.Run(port); err != nil {
		panic(err)
	}
	logServiceStart(appCtx)
	return appCtx
}

func logServiceStart(appCtx *ApplicationContext) {
	endTime := time.Now()
	elapsedTime := endTime.Sub(appCtx.StartTime)
	appCtx.logger.WithFields(logrus.Fields{
		"StartupDate":      appCtx.StartupDate,
		"ProjectName":      appCtx.DisplayName,
		"Profile":          appCtx.Profile,
		"GoVersion":        appCtx.GoVersion,
		"PID":              appCtx.Pid,
		"completedSeconds": fmt.Sprintf("%dm %ds", int(elapsedTime.Minutes()), int(elapsedTime.Seconds())%60),
	}).Info("Application Start!")
}
