package core

import (
	"fmt"
	"testing"
)

func TestFormatAntNumber(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{0, "00"},
		{1, "01"},
		{9, "09"},
		{10, "10"},
		{11, "11"},
		{99, "99"},
		{100, "100"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input_%d", tc.input), func(t *testing.T) {
			result := formatAntNumber(tc.input)
			if result != tc.expected {
				t.Errorf("formatAntNumber(%d) = %s; want %s", tc.input, result, tc.expected)
			}
		})
	}
}
