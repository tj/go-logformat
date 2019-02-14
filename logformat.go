// Package logformat provides some ad-hoc log formatting for some of my projects.
package logformat

import (
	"fmt"
	"sort"
	"time"

	"github.com/aybabtme/rgbterm"
)

// rgb struct.
type rgb struct {
	r uint8
	g uint8
	b uint8
}

// config is the formatter configuration.
type config struct {
	color rgb
}

// Option function.
type Option func(*config)

// WithColor option sets the primary color.
func WithColor(r, g, b uint8) Option {
	return func(v *config) {
		v.color = rgb{r, g, b}
	}
}

// Expanded returns a value in the expanded format.
func Expanded(v interface{}, options ...Option) string {
	c := config{
		color: rgb{112, 78, 251},
	}
	for _, o := range options {
		o(&c)
	}
	return format(v, "  ", &c)
}

// format returns a formatted value with prefix.
func format(v interface{}, prefix string, c *config) string {
	switch v := v.(type) {
	case map[string]interface{}:
		return formatMap(v, prefix, c)
	case []interface{}:
		return formatSlice(v, prefix, c)
	case time.Time:
		return v.Format(time.Stamp)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatMap returns a formatted map.
func formatMap(m map[string]interface{}, prefix string, c *config) string {
	s := ""
	keys := mapKeys(m)
	for _, k := range keys {
		v := m[k]
		k = rgbterm.FgString(k, c.color.r, c.color.g, c.color.b)
		if isComposite(v) {
			s += fmt.Sprintf("%s%s:\n%s", prefix, k, format(v, prefix+"  ", c))
		} else {
			s += fmt.Sprintf("%s%s: %s\n", prefix, k, format(v, prefix+"  ", c))
		}
	}
	return s
}

// formatSlice returns a formatted slice.
func formatSlice(v []interface{}, prefix string, c *config) string {
	s := ""
	for i, v := range v {
		if isComposite(v) {
			s += fmt.Sprintf("%s%d:\n%s", prefix, i, format(v, prefix+"  ", c))
		} else {
			s += fmt.Sprintf("%s- %v\n", prefix, v)
		}
	}
	return s
}

// isComposite returns true if the value is a composite.
func isComposite(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}:
		return true
	case []interface{}:
		return true
	default:
		return false
	}
}

// mapKeys returns map keys, sorted ascending.
func mapKeys(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}
