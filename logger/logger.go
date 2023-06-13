package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

const (
	LogLevelDebug = iota
	LogLevelRelease
)

var Logger *log.Logger
var LogLevel int

func init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	LogLevel = LogLevelRelease
	//LogLevel = LogLevelDebug
}

func Warnln(msg string) {
	Logger.Println(ColorYellow + msg + ColorReset)
}

func Fatalln(msg string) {
	Logger.Fatalln(ColorRed + msg + ColorReset)
}

func Infoln(msg string) {
	Logger.Println(ColorGreen + msg + ColorReset)
}

func Debugln(msg string) {
	if LogLevel != LogLevelDebug {
		return
	} else {
		Logger.Println(ColorBlue + msg + ColorReset)
	}
}

func WarnMsg(msg string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, msg, ColorReset)
}
