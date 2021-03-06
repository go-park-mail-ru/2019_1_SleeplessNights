package console

import (
	"fmt"
	"github.com/fatih/color"
)

var successColor        *color.Color
var errorColor          *color.Color
var messageColor        *color.Color
var titleColor          *color.Color
var predicateTrueColor  *color.Color
var predicateFalseColor *color.Color

func init() {
	successColor        = color.New(color.FgGreen).Add(color.Bold)
	errorColor          = color.New(color.FgRed).Add(color.Bold)
	messageColor        = color.New(color.FgBlue).Add(color.Bold)
	titleColor          = color.New(color.FgMagenta).Add(color.Bold).Add(color.Underline)
	predicateTrueColor  = color.New(color.FgCyan)
	predicateFalseColor = color.New(color.FgRed)
}

func Success(format string, a ...interface{}) {
	_, err := successColor.Printf(format + "\n", a...)
	if err != nil {
		fmt.Printf(format + "\n", a...)
	}
}

func Error(format string, a ...interface{}) {
	_, err := errorColor.Printf(format + "\n", a...)
	if err != nil {
		fmt.Printf(format + "\n", a...)
	}
}

func Message(format string, a ...interface{}) {
	_, err := messageColor.Printf(format + "\n", a...)
	if err != nil {
		fmt.Printf(format + "\n", a...)
	}
}

func Title(format string, a ...interface{}) {
	_, err := titleColor.Printf(format + "\n", a...)
	if err != nil {
		fmt.Printf(format + "\n", a...)
	}
}

func Predicate(key bool, format string, a ...interface{}) {
	var err error
	if key {
		_, err = predicateTrueColor.Printf("[ok] " + format+"\n", a...)
	} else {
		_, err = predicateFalseColor.Printf("[fail] " + format + "\n", a...)
	}
	if err != nil {
		fmt.Print(key, "")
		fmt.Printf(format + "\n", a...)
	}
}
