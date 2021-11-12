package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidUsageSegmentValidator(t *testing.T) {
	usageSegmentsData := []struct {
		testCaseDescriptor string
		testValue          int
		shouldBeValid bool
	} {
		{ "valid usage segment", 30, true },
		{ "invalid usage segment (negative)", -5, false },
		{ "invalid usage segment (>100)", 101, false },
	}

	for _, testCase := range usageSegmentsData {
		result := IsValidUsageSegment(testCase.testValue)

		assert.Equal(t, result, testCase.shouldBeValid,
			fmt.Sprintf("Unexpected usage segment validation result for case: %v", testCase.testCaseDescriptor))
	}
}
