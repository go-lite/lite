package lite

import (
	"testing"
)

func TestStatusMessage(t *testing.T) {
	tt := []struct {
		status   int
		expected string
	}{
		{
			status:   400,
			expected: "Bad Request",
		},
		{
			status:   401,
			expected: "Unauthorized",
		},
		{
			status:   404,
			expected: "Not Found",
		},
		{
			status:   409,
			expected: "Conflict",
		},
		{
			status:   500,
			expected: "Internal Server Error",
		},
		{
			status:   200,
			expected: "OK",
		},
		{
			status:   201,
			expected: "Created",
		},
		{
			status:   99,
			expected: "Unknown Status Code",
		},
		{
			status:   453,
			expected: "Unknown Status Code",
		},
	}

	for _, tc := range tt {
		t.Run(tc.expected, func(t *testing.T) {
			if got := StatusMessage(tc.status); got != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, got)
			}
		})
	}
}
