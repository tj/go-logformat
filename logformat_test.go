package logformat_test

import (
	"fmt"
	"testing"
	"time"

	logformat "github.com/tj/go-logformat"
)

var log = map[string]interface{}{
	// "timestamp": time.Now(),
	"timestamp": time.Now().Format(time.RFC3339),
	"message":   "response",
	"app":       "up-api",
	"version":   27,
	"level":     "info",
	"ip":        "35.190.145.206",
	"plugin":    "logs",
	"size":      "7998",
	"id":        "178f5348-304b-11e9-88af-9be94ad4eff3",
	"stage":     "production",
	"duration":  359,
	"cart": map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{
				"name": "Ferret food",
				"cost": 10.99,
			},
			map[string]interface{}{
				"name": "Cat food",
				"cost": 25.99,
			},
		},
		"total": 15.99,
		"paid":  false,
	},
	"pets": map[string]interface{}{
		"list": []interface{}{
			"Tobi",
			"Loki",
			"Jane",
		},
	},
	"method": "GET",
	"commit": "1d652f6",
	"path":   "/install",
	"query":  "",
	"region": "us-west-2",
	"status": "200",
}

// Test compact logs.
func TestCompact(t *testing.T) {
	fmt.Printf("%s\n\n", logformat.Compact(log))
}

// Test compact logs with flattening.
func TestCompact_WithFlatten(t *testing.T) {
	fmt.Printf("%s\n\n", logformat.Compact(log, logformat.WithFlatten(true)))
}

// Test expanded logs.
func TestExpanded(t *testing.T) {
	fmt.Printf("%s\n", logformat.Expanded(log))
}
