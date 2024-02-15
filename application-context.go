package goat

import "time"

type GoatApplication struct {
	startDateTime time.Time
}

func (app *GoatApplication) Run() {
	app.startDateTime = time.Now()
	// 기본 값 가져오기
	// propertySource 가져오기
	// RunListeners 가져오기
	// env 생성하기
	// context 생성하기
	// runListeners 실행하기(context, env)
	// context 반환하기
}

/**
여기서 알아야할 것
context 반환이 필요한지
context 실행이 필요한지
- 이건 스레드로 돌리는거도 아니다 프로그램이 정지되는 것을 막기 위해 잡고 있는 형태.
그렇기 때문에 context반환은 의미가 없다.
*/
