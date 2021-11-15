package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapUsageSegment(t *testing.T) {
	capUsageSegmentData := []struct {
		testCaseDescriptor string
		testValue          int
		expectedValue      int
	}{
		{"where capping not needed", 30, 30},
		{"negative, cap at 0", -5, 0},
		{"above 100, cap at 100", 101, 100},
	}

	for _, testCase := range capUsageSegmentData {
		result := CapUsageSegment(testCase.testValue)

		assert.Equal(t, result, testCase.expectedValue,
			fmt.Sprintf("Unexpected usage segment capping result: %v", testCase.testCaseDescriptor))
	}
}
