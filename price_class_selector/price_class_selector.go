package price_class_selector

import "time"

func GiveCorrectPrice(currentReadingUnixTimestamp int64, meteringPointTypeId int) float64 {
	if meteringPointTypeId == 2 {
		return GasPrice
	}

	readingDate := time.Unix(currentReadingUnixTimestamp, 0)

	// See if the reading was done during a weekday
	if int(readingDate.Weekday()) >= 1 && int(readingDate.Day()) < 7 {
		readingDaySevenAm := time.Date(readingDate.YearDay(), readingDate.Month(), readingDate.Day(),
			7, 0, 0, 0, time.Local)
		readingDayElevenPm := time.Date(readingDate.YearDay(), readingDate.Month(), readingDate.Day(),
			23, 0, 0, 0, time.Local)

		if readingDate.Before(readingDaySevenAm) && readingDate.After(readingDayElevenPm) {
			return ElectricPriceWkdayElevenPmTilSevenAm
		} else {
			return ElectricPriceWkdaySevenAmTilElevenPm
		}
	}

	return ElectricPriceWeekend
}
