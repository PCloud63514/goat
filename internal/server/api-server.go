package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goat/internal/application"
	"net/http"
)

const (
	DEFAULT_GIN_MODE    = gin.DebugMode
	DEFAULT_SERVER_ADDR = ":8080"
)

var (
	router *gin.Engine
	server *http.Server
)

func init() {
	addr := application.GetOrDefaultPropertyString("SERVER_ADDR", DEFAULT_SERVER_ADDR)
	mode := application.GetOrDefaultPropertyString("GIN_MODE", DEFAULT_GIN_MODE)

	gin.SetMode(mode)
	router = gin.Default()

	server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	application.AddStartEventCallback(listen)
	application.AddDestroyEventCallback(shutdown)
}

func listen() {
	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("listen: %s", err)
		}
	}()
}

func shutdown(ctx context.Context) {
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
