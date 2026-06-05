PROJECT_NAME := goat

verify:
	go mod tidy
	go test -v ./...

cover:
	# "패키지 레벨의 커버리지"와 "테스트 코드가 존재하는 파일 별 커버리지" 출력
	# "패키지 레벨의 커버리지"는 달성 목표 60% 이상, "파일 별 커버리지"는 달성 목표 60% 이상
	go run ./tools/coverreport -color -order=coverage