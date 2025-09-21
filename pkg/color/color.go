package color

import (
	"fmt"
)

// Simple color definitions for terminal output
var (
	Reset      = "\033[0m"
	RedCode    = "\033[31m"
	GreenCode  = "\033[32m"
	YellowCode = "\033[33m"
	BlueCode   = "\033[34m"
	PurpleCode = "\033[35m"
	CyanCode   = "\033[36m"
	WhiteCode  = "\033[37m"
)

// Color represents a color instance
type Color struct {
	code string
}

// Color constants
const (
	FgRed     = iota + 31
	FgGreen   = 32
	FgYellow  = 33
	FgBlue    = 34
	FgMagenta = 35
	FgCyan    = 36
	FgWhite   = 37
	FgHiCyan  = 96
)

// Attribute constants
const (
	Bold = 1
)

// New creates a new color with attributes
func New(attrs ...int) *Color {
	code := "\033["
	for i, attr := range attrs {
		if i > 0 {
			code += ";"
		}
		code += fmt.Sprintf("%d", attr)
	}
	code += "m"
	return &Color{code: code}
}

// Printf prints with color formatting
func (c *Color) Printf(format string, args ...interface{}) {
	fmt.Printf(c.code+format+Reset, args...)
}

// Println prints a line with color formatting
func (c *Color) Println(text string) {
	fmt.Println(c.code + text + Reset)
}

// Sprint returns a colored string
func (c *Color) Sprint(text string) string {
	return c.code + text + Reset
}

// Convenience functions
func Red(format string, args ...interface{}) {
	fmt.Printf(RedCode+format+Reset+"\n", args...)
}

func Green(format string, args ...interface{}) {
	fmt.Printf(GreenCode+format+Reset+"\n", args...)
}

func Yellow(format string, args ...interface{}) {
	fmt.Printf(YellowCode+format+Reset+"\n", args...)
}

func Blue(format string, args ...interface{}) {
	fmt.Printf(BlueCode+format+Reset+"\n", args...)
}

func Cyan(format string, args ...interface{}) {
	fmt.Printf(CyanCode+format+Reset+"\n", args...)
}

func White(format string, args ...interface{}) {
	fmt.Printf(WhiteCode+format+Reset+"\n", args...)
}
