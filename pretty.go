// Package pretty provides pretty printing for unmarshaled JSON Go types,
// thus it is very niche and probably not useful for you, I just need it for a few things.
package pretty

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aybabtme/rgbterm"
)

// Print a formatted value.
func Print(v interface{}) error {
	_, err := os.Stdout.WriteString(Format(v))
	return err
}

// Format returns a formatted value.
func Format(v interface{}) string {
	return format(v, "  ")
}

// format returns a formatted value with prefix.
func format(v interface{}, prefix string) string {
	switch v := v.(type) {
	case map[string]interface{}:
		return formatMap(v, prefix)
	case []interface{}:
		return formatSlice(v, prefix)
	case time.Time:
		return v.Format(time.Stamp)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatMap returns a formatted map.
func formatMap(m map[string]interface{}, prefix string) string {
	s := ""
	keys := mapKeys(m)
	for _, k := range keys {
		v := m[k]
		k = rgbterm.FgString(k, 112, 78, 251)
		if isComposite(v) {
			s += fmt.Sprintf("%s%s:\n%s", prefix, k, format(v, prefix+"  "))
		} else {
			s += fmt.Sprintf("%s%s: %s\n", prefix, k, format(v, prefix+"  "))
		}
	}
	return s
}

// formatSlice returns a formatted slice.
func formatSlice(v []interface{}, prefix string) string {
	s := ""
	for i, v := range v {
		if isComposite(v) {
			s += fmt.Sprintf("%s%d:\n%s", prefix, i, format(v, prefix+"  "))
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
