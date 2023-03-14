package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/api/quiz/app"
	"net/http"
)

func NewQuizApi(createQuizService app.CreateQuizService) QuizApi {
	return &quizApi{
		createQuizService: createQuizService,
	}
}

type QuizApi interface {
	RegisterHandler(router *gin.Engine)
	CreateQuiz(c *gin.Context)
	UpdateQuiz(c *gin.Context)
	DeleteQuiz(c *gin.Context)
	GetQuiz(c *gin.Context)
	GetQuizList(c *gin.Context)
}

type quizApi struct {
	createQuizService app.CreateQuizService
}

func (api *quizApi) RegisterHandler(router *gin.Engine) {
	router.POST("/api/quizzes", api.CreateQuiz)
	router.PUT("/api/quizzes/:id", api.UpdateQuiz)
	router.DELETE("/api/quizzes/:id", api.DeleteQuiz)
	router.GET("/api/quizzes/:id", api.GetQuiz)
	router.GET("/api/quizzes", api.GetQuizList)
}

func (api *quizApi) CreateQuiz(c *gin.Context) {
	var request app.CreateQuizRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = api.createQuizService.CreateQuiz(&request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}
func (api *quizApi) UpdateQuiz(c *gin.Context)  {}
func (api *quizApi) DeleteQuiz(c *gin.Context)  {}
func (api *quizApi) GetQuiz(c *gin.Context)     {}
func (api *quizApi) GetQuizList(c *gin.Context) {}
