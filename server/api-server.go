package server

import (
	"context"
	"errors"
	"github.com/PCloud63514/goat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	DEFAULT_GIN_MODE    = gin.DebugMode
	DEFAULT_SERVER_ADDR = ":9090"
)

var (
	router *gin.Engine
	server *http.Server
)

func init() {
	addr := goat.GetPropertyString("SERVER_ADDR", DEFAULT_SERVER_ADDR)
	mode := goat.GetPropertyString("GIN_MODE", DEFAULT_GIN_MODE)

	gin.SetMode(mode)
	router = gin.New()
	router.Use(CustomLogger(), gin.Recovery())
	server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	goat.AddHandlerFunc(listen, goat.HandlerType_Starting)
	goat.AddHandlerFunc(shutdown, goat.HandlerType_Stop)
}

func listen(ctx context.Context, env *goat.Environment) {
	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("listen msg: %s", err)
		}
	}()
}

func shutdown(ctx context.Context, env *goat.Environment) {
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown", err)
	}
}

func Router() *gin.Engine {
	return router
}

func Server() *http.Server {
	return server
}
