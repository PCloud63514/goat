package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

func (api *UserApi) RegisterHandler(router *gin.Engine) {
	router.POST("/api/users/sign-up", api.MemberShip)
	router.GET("/api/users/:id", api.GetUser)
}

func (api *UserApi) MemberShip(c *gin.Context) {

}

func (api *UserApi) GetUser(c *gin.Context) {
	id, _ := c.Params.Get("id")
	fmt.Println(id)
}
