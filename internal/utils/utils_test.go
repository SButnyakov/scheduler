package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLMDA(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected int
	}{
		{
			name:     "Empty array",
			input:    [][]int{},
			expected: 0,
		},
		{
			name:     "Single empty sub-array",
			input:    [][]int{{}},
			expected: 0,
		},
		{
			name:     "Single sub-array with elements",
			input:    [][]int{{1, 2, 3}},
			expected: 3,
		},
		{
			name:     "Multiple sub-arrays with elements",
			input:    [][]int{{1, 2, 3}, {4, 5}, {6}},
			expected: 6,
		},
		{
			name:     "Multiple sub-arrays with empty sub-arrays",
			input:    [][]int{{1, 2, 3}, {}, {4, 5}, {}, {6}},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LMDA(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsChannelClosed(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() chan struct{}
		expected bool
	}{
		{
			name: "Open channel",
			setup: func() chan struct{} {
				return make(chan struct{})
			},
			expected: false,
		},
		{
			name: "Closed channel",
			setup: func() chan struct{} {
				ch := make(chan struct{})
				close(ch)
				return ch
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.setup()
			result := IsChannelClosed(ch)
			assert.Equal(t, tt.expected, result)
		})
	}
}
