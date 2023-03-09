package pkg

import "github.com/gin-gonic/gin"

type Api interface {
	RegisterHandler(router *gin.Engine)
}
