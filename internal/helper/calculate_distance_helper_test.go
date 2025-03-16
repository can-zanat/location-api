package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{"Same Point", 0, 0, 0, 0, 0},
		{"Istanbul to Ankara", 41.0082, 28.9784, 39.9334, 32.8597, 351.0},
		{"New York to Los Angeles", 40.7128, -74.0060, 34.0522, -118.2437, 3935.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)

			assert.InDelta(t, tt.expected, result, 5.0) // ±5 km tolerans
		})
	}
}
