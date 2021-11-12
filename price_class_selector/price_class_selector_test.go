package price_class_selector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGiveCorrectPrice(t *testing.T) {
	priceClassReadingData := []struct {
		testCaseDescriptor string
		unixTimestamp      int64
		meteringPointId int
		expectedPrice float64
		shouldBeValid bool
	} {
		{
			"electric price for weekday after 7 AM and before 11 PM w/correct price class",
			1415964600,
			1,
			ElectricPriceWkdaySevenAmTilElevenPm,
			true,
		},
		{
			"electric price for weekday after 7 AM and before 11 PM w/incorrect price class",
			1636735707,
			1,
			ElectricPriceWkdayElevenPmTilSevenAm,
			false,
		},
		{
			"electric price for weekend w/correct price class",
			1636822107,
			1,
			ElectricPriceWeekend,
			true,
		},
		{
			"gas price w/correct price class",
			1636822107,
			2,
			GasPrice,
			true,
		},
		{
			"gas price w/incorrect price class",
			1636822107,
			2,
			ElectricPriceWkdaySevenAmTilElevenPm,
			false,
		},
	}

	for _, testCase := range priceClassReadingData {
		result := GiveCorrectPrice(testCase.unixTimestamp, testCase.meteringPointId) == testCase.expectedPrice

		assert.Equal(t, result, testCase.shouldBeValid,
			fmt.Sprintf("Unexpected price class giver result for case: %v", testCase.testCaseDescriptor))
	}
}
