package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"main/api/quiz/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCreateQuizService struct {
	mock.Mock
}

func (m *mockCreateQuizService) CreateQuiz(request *app.CreateQuizRequest) (*app.CreateQuizResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*app.CreateQuizResponse), args.Error(1)
}

func TestNewQuizApi(t *testing.T) {
	suite.Run(t, new(QuizApiTestSuite))
}

type QuizApiTestSuite struct {
	suite.Suite
	api                   QuizApi
	router                *gin.Engine
	mockCreateQuizService *mockCreateQuizService
}

func (suite *QuizApiTestSuite) SetupTest() {
	router := gin.Default()
	_mockCreateQuizService := new(mockCreateQuizService)
	suite.mockCreateQuizService = _mockCreateQuizService
	suite.api = NewQuizApi(_mockCreateQuizService)
	suite.api.RegisterHandler(router)
	suite.router = router
}

// [퀴즈 생성] 요청이 성공했을 경우 Status Created(201)을 반환합니다.
func (suite *QuizApiTestSuite) TestAddQuiz_return_CreatedStatus() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/quizzes", nil)
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

// [퀴즈 생성] 요청문은 CreateQuizService에 전달되어야한다.
func (suite *QuizApiTestSuite) TestAddQuiz_passes_request_CreateQuizService() {
	givenRequest := &app.CreateQuizRequest{
		Title:   "테스트",
		Desc:    "테스트 설명",
		Timer:   0,
		Answer:  1,
		Select1: "고양이",
		Select2: "강아지",
		Select3: "도마뱀",
		Select4: "새",
	}
	suite.mockCreateQuizService.On("CreateQuiz", mock.Anything).Return(nil)

	reqBody, _ := json.Marshal(givenRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/quizzes", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	arg := suite.mockCreateQuizService.Calls[0].Arguments.Get(0)
	_, ok := arg.(*app.CreateQuizRequest)
	assert.True(suite.T(), ok)
}
