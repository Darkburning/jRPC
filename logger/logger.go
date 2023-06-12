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
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func Warnln(msg string) {
	Logger.Println(ColorYellow + msg + ColorReset)
}

func Fatalln(msg string) {
	Logger.Println(ColorRed + msg + ColorReset)
}

func Infoln(msg string) {
	Logger.Println(ColorGreen + msg + ColorReset)
}

func WarnMsg(msg string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, msg, ColorReset)
}
