package application

import (
	"flag"
	"time"
)

/*
+-------------+
|  Arguments  |
+-------------+
*/

const (
	FLAG_PROFILE_NAME  = "profile"
	FLAG_PROFILE_VALUE = "default"
	FLAG_PROFILE_USAGE = "프로필 속성입니다. 기본 값은 default 입니다."
)

type applicationArguments struct {
	Profile     string
	ProjectName string
	options     map[string]string
	StartTime   time.Time
}

func newApplicationArguments(startTime time.Time) *applicationArguments {
	profileFlag := flag.String(FLAG_PROFILE_NAME, FLAG_PROFILE_VALUE, "프로필")
	flag.Parse()
	_options := make(map[string]string)
	flag.Visit(func(f *flag.Flag) {
		_options[f.Name] = f.Value.String()
	})

	return &applicationArguments{
		Profile:   *profileFlag,
		options:   _options,
		StartTime: startTime,
	}
}
