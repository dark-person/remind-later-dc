package timeparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	// Normal test
	type testCase struct {
		input  string
		hour   int
		minute int
		second int
	}

	tests := []testCase{
		{"1h2m3s", 1, 2, 3},
		{"2m3s", 0, 2, 3},
		{"1h2m", 1, 2, 0},
		{"1h", 1, 0, 0},
		{"2m", 0, 2, 0},
		{"3s", 0, 0, 3},
		{"1h3s", 1, 0, 3},
	}

	for _, tt := range tests {
		hour, minute, second, ok := ExtractDuration(tt.input)

		assert.EqualValues(t, tt.hour, hour)
		assert.EqualValues(t, tt.minute, minute)
		assert.EqualValues(t, tt.second, second)
		assert.True(t, ok)
	}

	// Case expected to failed to parse
	errTests := []string{"", "abc", "3s2m"}
	for _, tt := range errTests {
		hour, minute, second, ok := ExtractDuration(tt)

		assert.EqualValues(t, 0, hour)
		assert.EqualValues(t, 0, minute)
		assert.EqualValues(t, 0, second)
		assert.False(t, ok)
	}
}
