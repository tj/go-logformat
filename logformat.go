// Package logformat provides some ad-hoc log formatting for some of my projects.
package logformat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nqd/flat"
	"github.com/tj/go-logformat/internal/colors"
)

// FormatFunc is a function for formatting values.
type FormatFunc func(string) string

// Formatters is a map of formatting functions.
type Formatters map[string]FormatFunc

// DefaultFormatters is a set of default formatters.
var DefaultFormatters = Formatters{
	// Levels.
	"debug":     colors.Gray,
	"info":      colors.Purple,
	"warn":      colors.Yellow,
	"warning":   colors.Yellow,
	"error":     colors.Red,
	"fatal":     colors.Red,
	"critical":  colors.Red,
	"emergency": colors.Red,

	// Values.
	"string": colors.None,
	"number": colors.None,
	"bool":   colors.None,
	"date":   colors.Gray,

	// Fields.
	"object.key":       colors.Purple,
	"object.separator": colors.Gray,
	"object.value":     colors.None,

	// Arrays.
	"array.delimiter": colors.Gray,
	"array.separator": colors.Gray,

	// Special fields.
	"message": colors.None,
	"program": colors.Gray,
	"stage":   colors.Gray,
	"version": colors.Gray,
}

// NoColor is a set of formatters resulting in no colors.
var NoColor = Formatters{
	// Levels.
	"debug":     colors.None,
	"info":      colors.None,
	"warn":      colors.None,
	"warning":   colors.None,
	"error":     colors.None,
	"fatal":     colors.None,
	"critical":  colors.None,
	"emergency": colors.None,

	// Values.
	"string": colors.None,
	"number": colors.None,
	"bool":   colors.None,
	"date":   colors.None,

	// Fields.
	"object.key":       colors.None,
	"object.separator": colors.None,
	"object.value":     colors.None,

	// Arrays.
	"array.delimiter": colors.None,
	"array.separator": colors.None,

	// Special fields.
	"message": colors.None,
	"program": colors.None,
	"stage":   colors.None,
	"version": colors.None,
}

// config is the formatter configuration.
type config struct {
	format  Formatters
	flatten bool
}

// Option function.
type Option func(*config)

// WithFormatters overrides the default formatters.
func WithFormatters(v Formatters) Option {
	return func(c *config) {
		c.format = v
	}
}

// WithFlatten toggles flattening of fields.
func WithFlatten(v bool) Option {
	return func(c *config) {
		c.flatten = v
	}
}

// newConfig returns config with options applied.
func newConfig(options ...Option) *config {
	c := &config{
		format: DefaultFormatters,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// Compact returns a value in the compact format.
func Compact(m map[string]interface{}, options ...Option) string {
	return compact(m, newConfig(options...))
}

// compact returns a formatted value.
func compact(v interface{}, c *config) string {
	switch v := v.(type) {
	case map[string]interface{}:
		return compactMap(v, c)
	case []interface{}:
		return compactSlice(v, c)
	default:
		return primitive(v, c)
	}
}

// Prefix returns a prefix for log line special-casing, and removes those fields from the map.
func Prefix(m map[string]interface{}, options ...Option) string {
	s := ""
	c := newConfig(options...)

	// timestamp
	if v, ok := m["timestamp"].(string); ok {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			s += primitive(t, c) + " "
			delete(m, "timestamp")
		}
	}

	// level
	if v, ok := m["level"].(string); ok {
		format := c.format[v]
		if format != nil {
			s += bold(format(strings.ToUpper(v[:4]))) + " "
			delete(m, "level")
		}
	}

	// application
	if v, ok := m["application"].(string); ok {
		s += c.format["application"](v) + " "
		delete(m, "application")
	}

	// stage
	if v, ok := m["stage"].(string); ok {
		s += c.format["stage"](v) + " "
		delete(m, "stage")
	}

	// version
	if v, ok := m["version"]; ok {
		switch v := v.(type) {
		case string:
			s += c.format["version"](v) + " "
		case float64:
			s += c.format["version"](strconv.FormatFloat(v, 'f', -1, 64)) + " "
		}
		delete(m, "version")
	}

	// message
	if v, ok := m["message"].(string); ok {
		s += c.format["message"](v)
		delete(m, "message")
	}

	return s
}

// compactMap returns a formatted map.
func compactMap(m map[string]interface{}, c *config) string {
	m = maybeFlatten(m, c)
	s := ""
	keys := mapKeys(m)
	for i, k := range keys {
		v := m[k]
		s += c.format["object.key"](k)
		s += c.format["object.separator"]("=")
		if isComposite(v) {
			s += c.format["object.value"](compact(v, c))
		} else {
			s += compact(v, c)
		}
		if i < len(keys)-1 {
			s += " "
		}
	}
	return s
}

// compactSlice returns a formatted slice.
func compactSlice(v []interface{}, c *config) string {
	s := c.format["array.delimiter"]("[")
	for i, v := range v {
		if i > 0 {
			s += c.format["array.separator"](", ")
		}
		s += compact(v, c)
	}
	return s + c.format["array.delimiter"]("]")
}

// Expanded returns a value in the expanded format.
func Expanded(m map[string]interface{}, options ...Option) string {
	return expanded(m, "  ", newConfig(options...))
}

// expanded returns a formatted value with prefix.
func expanded(v interface{}, prefix string, c *config) string {
	switch v := v.(type) {
	case map[string]interface{}:
		return expandedMap(v, prefix, c)
	case []interface{}:
		return expandedSlice(v, prefix, c)
	default:
		return primitive(v, c)
	}
}

// expandedMap returns a formatted map.
func expandedMap(m map[string]interface{}, prefix string, c *config) string {
	m = maybeFlatten(m, c)
	s := ""
	keys := mapKeys(m)
	for _, k := range keys {
		v := m[k]
		k = c.format["object.key"](k)
		d := c.format["object.separator"](":")
		if isComposite(v) {
			s += fmt.Sprintf("%s%s%s\n%s", prefix, k, d, expanded(v, prefix+"  ", c))
		} else {
			s += fmt.Sprintf("%s%s%s %s\n", prefix, k, d, expanded(v, prefix+"  ", c))
		}
	}
	return s
}

// expandedSlice returns a formatted slice.
func expandedSlice(v []interface{}, prefix string, c *config) string {
	s := ""
	for _, v := range v {
		d := c.format["array.separator"]("-")
		if isComposite(v) {
			s += fmt.Sprintf("%s%s\n%s", prefix, d, expanded(v, prefix+"  ", c))
		} else {
			s += fmt.Sprintf("%s%s %v\n", prefix, d, primitive(v, c))
		}
	}
	return s
}

// primitive returns a formatted value.
func primitive(v interface{}, c *config) string {
	switch v := v.(type) {
	case string:
		if strings.ContainsAny(v, " \n\t") || strings.TrimSpace(v) == "" {
			return c.format["string"](strconv.Quote(v))
		} else {
			return c.format["string"](v)
		}
	case time.Time:
		return c.format["date"](formatDate(v))
	case bool:
		return c.format["bool"](strconv.FormatBool(v))
	case float64:
		return c.format["number"](strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return fmt.Sprintf("%v", v)
	}
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

// formatDate formats t relative to now.
func formatDate(t time.Time) string {
	return t.Format(`Jan 2` + dateSuffix(t) + ` 03:04:05pm`)
}

// bold string.
func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

// dateSuffix returns the date suffix for t.
func dateSuffix(t time.Time) string {
	switch t.Day() {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}

// maybeFlatten returns a the original or flattened map when configured to do so.
func maybeFlatten(m map[string]interface{}, c *config) map[string]interface{} {
	if c.flatten {
		m, _ = flat.Flatten(m, &flat.Options{Safe: true, Delimiter: "."})
	}
	return m
}
