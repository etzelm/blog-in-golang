package main

import (
	"testing" // Import the standard Go testing package
)

// TestRandRange checks if the randRange function generates a number
// within the specified min and max boundaries.
func TestRandRange(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name string
		min  int
		max  int
	}{
		{"PositiveRange", 1, 10},
		{"SingleValueRange", 5, 5},
		{"NegativeRange", -10, -1},
		{"MixedRange", -5, 5},
		{"ZeroMinRange", 0, 10},
		{"ZeroMaxRange", -10, 0},
	}

	// Iterate over each test case
	for _, tc := range testCases {
		// t.Run creates a sub-test, making it easier to identify
		// which specific case failed if there are multiple failures.
		t.Run(tc.name, func(t *testing.T) {
			// Call the function we want to test
			result := randRange(tc.min, tc.max)

			// Check if the result is less than the minimum allowed value
			if result < tc.min {
				t.Errorf("randRange(%d, %d) = %d; want value >= %d", tc.min, tc.max, result, tc.min)
			}

			// Check if the result is greater than the maximum allowed value
			if result > tc.max {
				t.Errorf("randRange(%d, %d) = %d; want value <= %d", tc.min, tc.max, result, tc.max)
			}
		})
	}
}
