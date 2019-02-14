// Package colors provides some colors :)
package colors

import (
	color "github.com/aybabtme/rgbterm"
)

// None string.
func None(s string) string {
	return s
}

// Gray string.
func Gray(s string) string {
	return color.FgString(s, 150, 150, 150)
}

// Blue string.
func Blue(s string) string {
	return color.FgString(s, 32, 115, 191)
}

// Cyan string.
func Cyan(s string) string {
	return color.FgString(s, 25, 133, 152)
}

// Green string.
func Green(s string) string {
	return color.FgString(s, 48, 137, 65)
}

// Red string.
func Red(s string) string {
	return color.FgString(s, 194, 37, 92)
}

// Yellow string.
func Yellow(s string) string {
	return color.FgString(s, 252, 196, 25)
}

// Purple string.
func Purple(s string) string {
	return color.FgString(s, 96, 97, 190)
}
