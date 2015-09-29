package io

import (
	"fmt"
	"os"
)

func Info(args ...interface{}) {
	fmt.Print("\033[1m-----> ")
	args = append(args, "\033[0m")
	fmt.Println(args...)
}

func Infof(format string, args ...interface{}) {
	fmt.Print("\033[1m-----> ")
	fmt.Printf(format+"\033[0m", args...)
}

func Print(args ...interface{}) {
	fmt.Print("       ")
	fmt.Println(args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Print("       ")
	fmt.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	fmt.Print(" !     ")
	fmt.Printf(format, args...)
}

func Error(args ...interface{}) {
	fmt.Print(" !     ")
	fmt.Println(args...)
	os.Exit(1)
}
