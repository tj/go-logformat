package logformat_test

import (
	"testing"
	"time"

	logformat "github.com/tj/go-logformat"
)

var object = map[string]interface{}{
	"timestamp": time.Now(),
	"message":   "response",
	"app":       "up-api",
	"version":   27,
	"fields": map[string]interface{}{
		"ip":       "35.190.145.206",
		"plugin":   "logs",
		"size":     "7998",
		"id":       "178f5348-304b-11e9-88af-9be94ad4eff3",
		"stage":    "production",
		"duration": 359,
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
	},
}

// Test formatting.
func TestPrint(t *testing.T) {
	logformat.Print(object)
}

// Test formatting with custom color.
func TestPrint_WithColor(t *testing.T) {
	logformat.Print(object, logformat.WithColor(164, 33, 78))
}
