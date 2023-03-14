package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewQuizApi(t *testing.T) {
	suite.Run(t, new(QuizApiTestSuite))
}

type QuizApiTestSuite struct {
	suite.Suite
	api    QuizApi
	router *gin.Engine
}

func (suite *QuizApiTestSuite) SetupTest() {
	router := gin.Default()
	suite.api = NewQuizApi()
	suite.api.RegisterHandler(router)
	suite.router = router
}

// [퀴즈 생성] 요청이 성공했을 경우 Status Created(201)을 반환합니다.
func (suite *QuizApiTestSuite) TestAddQuiz_returnCreatedStatus() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/quizzes", nil)
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}
