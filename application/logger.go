package application

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return logger
}

func ColorString(msg string, color Color) string {
	return fmt.Sprintf("%v%s%v", color, msg, COLOR_BACKGROUND)
}

type Color string

const (
	COLOR_BACKGROUND     Color = "\033[0m"
	COLOR_BLACK          Color = "\033[30m"
	COLOR_RED            Color = "\033[31m"
	COLOR_GREEN          Color = "\033[32m"
	COLOR_YELLOW         Color = "\033[33m"
	COLOR_BLUE           Color = "\033[34m"
	COLOR_MAGENTA        Color = "\033[35m"
	COLOR_CYAN           Color = "\033[36m"
	COLOR_WHITE          Color = "\033[37m"
	COLOR_BRIGHT_BLACK   Color = "\033[90m"
	COLOR_BRIGHT_RED     Color = "\033[91m"
	COLOR_BRIGHT_GREEN   Color = "\033[92m"
	COLOR_BRIGHT_YELLOW  Color = "\033[93m"
	COLOR_BRIGHT_BLUE    Color = "\033[94m"
	COLOR_BRIGHT_MAGENTA Color = "\033[95m"
	COLOR_BRIGHT_CYAN    Color = "\033[96m"
	COLOR_BRIGHT_WHITE   Color = "\033[97m"
	COLOR_DEFAULT        Color = "\033[39m"
	COLOR_BOLD           Color = "\033[1m"
	COLOR_UNDERLINE      Color = "\033[4m"
	COLOR_INVERSE        Color = "\033[7m"
	COLOR_RESET          Color = "\033[0m"
	COLOR_GRAY           Color = "\033[90m"
)
