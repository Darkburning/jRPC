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

var logger *log.Logger
var logLevel int

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logLevel = LogLevelRelease
	//logLevel = LogLevelDebug
}

func Warnln(msg string) {
	logger.Println(ColorYellow + msg + ColorReset)
}

func Fatalln(msg string) {
	logger.Fatalln(ColorRed + msg + ColorReset)
}

func Infoln(msg string) {
	logger.Println(ColorGreen + msg + ColorReset)
}

func Debugln(msg string) {
	if logLevel != LogLevelDebug {
		return
	} else {
		logger.Println(ColorBlue + msg + ColorReset)
	}
}

func WarnMsg(msg string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, msg, ColorReset)
}
