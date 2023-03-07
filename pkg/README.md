# pkg

외부 애플리케이션에서 사용되어도 괜찮은 라이브러리 코드입니다 (e.g., /pkg/mypubliclib). 다른 프로젝트는 이 라이브러리들이 작동할거라고 예상하고 임포트 할 것 이므로, 여기에 무언가를 넣기 전에 두번 고민하세요 :-) internal디렉터리는 개인적인 패키지들이 임포트 불가능하도록 하는 더 좋은 방법인데, 이유는 Go가 이를 강제하기 떄문입니다. /pkg 디렉터리는 그 디렉터리 안의 코드가 다른 사람들에 의해 사용되어도 안전하다고 명시적으로 보여주는 좋은 방법입니다. Travis Jeffery의 I'll take pkg over internal 블로그 포스트는 pkg 와 internal 디렉터리와 언제 쓰는게 맞을지에 대해 좋은 개요를 제공합니다.

또한 루트 디렉터리에 많은 Go가 아닌 컴포넌트와 디렉터리를 포함하고 있다면 Go 코드를 한 곳에 모아서 다양한 Go 툴들을 쉽게 실행할 수 있습니다 (이 발표들에서 언급되었던것 처럼: GopherCon EU 2018의 Best Practices for Industrial Programming, GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps 와 GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go).

이 프로젝트 레이아웃 패턴을 사용하는 유명한 Go 레포지토리들을 보고 싶다면 /pkg를 보세요. 이는 흔한 레이아웃 패턴이나, 보편적으로 받아드려지는 것은 아니며 Go 커뮤니티의 일부는 이를 추천하지 않습니다.

앱 프로젝트가 정말 작고 추가적인 레벨의 중첩이 많이 효과적이지 않다면 사용하지 않아도 괜찮습니다 (정말로 원하지 않는 한 :-)). 프로젝트가 더 커지고 루트 디렉터리가 꽤 바빠질 때 고려해보세요 (특히 Go가 아닌 앱 컴포넌트가 많이 있다면).


