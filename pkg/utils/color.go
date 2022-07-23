package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func MsgInfo(msg string) string {
	c := color.New(color.FgGreen, color.Bold)
	return c.Sprint(msg)
}

func MsgError(msg string) string {
	c := color.New(color.FgRed, color.Bold)
	return c.Sprint(msg)
}

func MsgWarning(msg string) string {
	c := color.New(color.FgYellow, color.Bold)
	return c.Sprint(msg)
}

func Bold(msg string) string {
	c := color.New(color.Bold)
	return c.Sprint(msg)
}

func PrintInfo(format string, msg ...interface{}) {
	c := color.New(color.FgGreen, color.Bold)
	m := fmt.Sprintf(format, msg...)
	c.Printf("%s %s", c.Sprint("[+]"), m)
}

func PrintError(format string, msg ...interface{}) {
	c := color.New(color.FgRed, color.Bold)
	m := fmt.Sprintf(format, msg...)
	c.Printf("%s %s", c.Sprint("[!]"), m)
}

func PrintWarning(format string, msg ...interface{}) {
	c := color.New(color.FgYellow, color.Bold)
	m := fmt.Sprintf(format, msg...)
	c.Printf("%s %s", c.Sprint("[*]"), m)
}
