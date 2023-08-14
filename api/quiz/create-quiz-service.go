package quiz

import "time"


type CreateQuizRequest struct {
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Timer    int64     `json:"timer"`
	Answer   int64     `json:"answer"`
	Select1  string    `json:"select1"`
	Select2  string    `json:"select2"`
	Select3  string    `json:"select3"`
	Select4  string    `json:"select4"`
	OpenAt   time.Time `json:"openAt"`
	OpenTime time.Time `json:"openTime"`
	EndTime  time.Time `json:"endTime"`
}
type CreateQuizResponse struct{}

func CreateQuiz(request *CreateQuizRequest) (*CreateQuizResponse, error) {

	return nil, nil
}
