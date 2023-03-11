package main

import (
	"main/application"
)

func main() {
	application.Start()
}

/*
Run 동작 시 목표
- arguments 불러오기
- environment 불러오기
- DB 객체 생성 및 연결하기
- 프로그램 동작에 필요한 인스턴스들 생성 및 등록하기
- 배치 스케줄러 객체 생성하기
- 배치 스케줄러에 잡 등록하기
- router에 RouteHandler 등록하기
- 기본 포트로 웹 서비스 실행 (router run)
- 배치 스케줄러 실행하기
*/
