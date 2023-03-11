package application

import (
	"os"
	"path/filepath"
)

/*
+-----------+
|   Utils   |
+-----------+
*/

func GetProjectName() string {
	// 실행 파일의 경로와 이름을 가져옵니다.
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 실행 파일의 디렉토리 경로를 가져옵니다.
	executableDir := filepath.Dir(executablePath)
	// 실행 파일의 디렉토리 이름을 가져옵니다.
	projectName := filepath.Base(executableDir)
	return projectName
}
