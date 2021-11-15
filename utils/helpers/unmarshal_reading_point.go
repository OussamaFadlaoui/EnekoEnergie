package helpers

import (
	"fmt"
	"github.com/OussamaFadlaoui/EnekoEnergie/types"
	"os"
	"strconv"
)


func UnmarshalReadingPoint(rawLine []string) types.ReadingPoint {
	errorMsgParsing := "[Error] Could not unmarshal reading point into type ReadingPoint"

	if len(rawLine) != 4 {
		fmt.Println("[Error] Invalid amount of columns read in a row of the inputs file")
		os.Exit(1)
	}

	meteringPointId, err := strconv.Atoi(rawLine[0])
	Check(err, errorMsgParsing)

	meteringTypeId, err := strconv.Atoi(rawLine[1])
	Check(err, errorMsgParsing)

	readingValue, err := strconv.Atoi(rawLine[2])
	Check(err, errorMsgParsing)

	createdAt, err := strconv.ParseInt(rawLine[3], 10, 64)
	Check(err, errorMsgParsing)

	return types.ReadingPoint{
		MeteringPointId: meteringPointId,
		MeteringTypeId:  meteringTypeId,
		ReadingValue:    readingValue,
		CreatedAt:       createdAt,
	}
}
