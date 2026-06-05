package main

import (
	"time"

	"github.com/PCloud63514/goat/environment"
	"github.com/PCloud63514/goat/logger"
	"github.com/PCloud63514/goat/profile"
)

//////////////////
// Goat Example //
//////////////////

var (
	startUpDateTime = time.Now()
	configurations  = []interface{}{
		&logger.Config{},
	}
)

func main() {
	environment.Reset(environment.Option{
		ResPath:        "res",
		Profiles:       profile.Get(),
		Configurations: configurations,
	})

	logger.Infof("Start Up Date Time: %s", startUpDateTime.Format(time.RFC3339))
	logger.Infof("Log Level: %s", environment.GetProperty("logger.level", "NONE"))
}
