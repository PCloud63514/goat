# internal

개인적인 애플리케이션과 라이브러리 코드. 다른 사람들이 애플리케이션이나 라이브러리에서 임포트 하기를 원하지 않는 코드들입니다. 이 레이아웃 패턴은 Go 컴파일러 자체에 강제됩니다. 더 알아보고 싶다면 Go 1.4 release notes 를 참고하세요. 이는 최상단 internal 디렉터리에 국한되는 것은 아닙니다. 프로젝트 트리의 모든 레벨에서 하나 이상의 internal 디렉터리를 가질 수 있습니다.

Internal 패키지에서 공유되는 내부 코드와 공유되지 않는 내부 코드를 구분하기 위해 부가적인 구조체를 추가할 수 있습니다. 필수사항은 아니지만(특히 작은 프로젝트에서는), 의도된 패키지 사용법을 보여주는 시각적인 단서를 남기는 것이 좋습니다. 실제 에플리케이션 코드는 /internal/app 디렉터리에 넣고 (e.g., /internal/app/myapp) , 그 앱들에서 공유되는 코드들은 /internal/pkg 디렉터리 (e.g., /internal/pkg/myprivlib) 에 넣을 수 있습니다.


## 참조
Project Layout: https://github.com/golang-standards/project-layout/blob/master/README_ko.md