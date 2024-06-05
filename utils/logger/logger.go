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

// fmt.Println("\033[31mThis is red text\033[0m")
// fmt.Println("\033[32mThis is green text\033[0m")
// fmt.Println("\033[33mThis is yellow text\033[0m")
// fmt.Println("\033[34mThis is blue text\033[0m")
// fmt.Println("\033[35mThis is purple text\033[0m")
// fmt.Println("\033[36mThis is cyan text\033[0m")
// fmt.Println("\033[1mThis is bold text\033[0m")
// fmt.Println("This is normal text")
