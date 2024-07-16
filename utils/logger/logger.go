package logger

import (
	"fmt"
	"os"
)

const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func colorPrint(color string, format string, a ...any) {
	coloredFormat := fmt.Sprintf("%s%s%s", color, format, reset)
	fmt.Printf(coloredFormat, a...)
}

func Success(format string, a ...any) {
	colorPrint(green, format, a...)
}

func Changed(format string, a ...any) {
	colorPrint(yellow, format, a...)
}

func Failed(format string, a ...any) {
	colorPrint(red, format, a...)
}

func Fatal(format string, a ...any) {
	colorPrint(red, format, a...)
	os.Exit(1)
}

func Ignore(format string, a ...any) {
	colorPrint(yellow, format, a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Println(a ...any) {
	fmt.Println(a...)
}
