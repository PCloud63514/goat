package quiz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goat/internal/server"
	"net/http"
)

func init() {
	server.Router().POST("/api/quizzes", AddQuiz)
	server.Router().PUT("/api/quizzes/:id", UpdateQuiz)
	server.Router().DELETE("/api/quizzes/:id", DeleteQuiz)
	server.Router().GET("/api/quizzes/:id", GetQuiz)
	server.Router().GET("/api/quizzes", GetQuizList)
}

func AddQuiz(c *gin.Context) {
	var request CreateQuizRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = CreateQuiz(&request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}
func UpdateQuiz(c *gin.Context)  {}
func DeleteQuiz(c *gin.Context)  {}
func GetQuiz(c *gin.Context)     {}
func GetQuizList(c *gin.Context) {}
