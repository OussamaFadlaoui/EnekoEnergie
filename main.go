package main

import (
	"encoding/csv"
	"fmt"
	"github.com/OussamaFadlaoui/EnekoEnergie/usage_validator"
	"github.com/OussamaFadlaoui/EnekoEnergie/utils"
	"github.com/OussamaFadlaoui/EnekoEnergie/utils/helpers"
	"io"
	"os"
)

func main() {
	inputFile, err := os.Open(ReadingInputsFilePath)

	helpers.Check(err, "[Error] Could not open inputs file for some reason")

	var reader = csv.NewReader(inputFile)
	var lineCount = 0
	var lastReadingPoint utils.ReadingPoint
	var currentReading utils.ReadingPoint
	var nextReading utils.ReadingPoint
	var calculatedUsageSegment int
	var usageSegments = make(map[int][]int)
	var invalidUsageSegmentIndices = make(map[int][]int)
	var flagNextReadingWillBeInvalid bool
	//var results map[int]float64

	for {
		if currentReading == (utils.ReadingPoint{}) {
			currentReadingRaw, err := reader.Read()
			helpers.Check(err, ErrorMsgReadingFile)
			currentReading = helpers.UnmarshalReadingPoint(currentReadingRaw)
		} else {
			currentReading = lastReadingPoint
		}

		nextReadingRaw, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(ErrorMsgReadingFile)
			os.Exit(1)
		} else if len(nextReadingRaw) != 4 {
			fmt.Println("[Error] Invalid amount of columns read in a row of the inputs file")
			os.Exit(1)
		}

		nextReading = helpers.UnmarshalReadingPoint(nextReadingRaw)

		// Flag found for this reading being invalid due to generating an invalid
		// usage segment value
		if flagNextReadingWillBeInvalid {
			lastReadingPoint = nextReading
			flagNextReadingWillBeInvalid = false
			lineCount++
			continue
		}

		if nextReading.MeteringPointId != currentReading.MeteringPointId {
			lastReadingPoint = nextReading
			continue
		}

		calculatedUsageSegment = nextReading.ReadingValue - currentReading.ReadingValue
		isValidUsageSegment := usage_validator.IsValidUsageSegment(calculatedUsageSegment)
		pointId := currentReading.MeteringPointId

		if !isValidUsageSegment {
			// Edge cases:
			if len(usageSegments[pointId]) == 0 {
				// First usage segment found was invalid
				usageSegments[pointId] = append(usageSegments[pointId], -1, -1)
				invalidUsageSegmentIndices[pointId] = append(invalidUsageSegmentIndices[pointId], 0, 1)
				flagNextReadingWillBeInvalid = true
			} else if len(usageSegments[pointId]) >= 2 && len(invalidUsageSegmentIndices[pointId]) == 0 {
				// Found invalid usage segment where we have two valid usage segments before it
				diffBetweenTwoLastSegments := usageSegments[pointId][1] - usageSegments[pointId][0]
				linearizedUsageSegment := usageSegments[pointId][1] + diffBetweenTwoLastSegments
				linearizedUsageSegment = helpers.CapUsageSegment(linearizedUsageSegment)

				usageSegments[pointId] = append(usageSegments[pointId], linearizedUsageSegment)
			} else if len(usageSegments[pointId]) == 1 && len(invalidUsageSegmentIndices[pointId]) == 0 {
				// Found invalid usage segment where there is only 1 valid usage segment before it
				firstUsageSegment := usageSegments[pointId][0]
				linearizedUsageSegment := helpers.CapUsageSegment(firstUsageSegment * 2)

				usageSegments[pointId] = append(usageSegments[pointId], linearizedUsageSegment)
			}
		} else {
			usageSegments[pointId] = append(usageSegments[pointId], calculatedUsageSegment)
		}

		lastReadingPoint = nextReading
		lineCount++
	}

	fmt.Println(usageSegments)

	fmt.Println("Finished reading file")
	err = inputFile.Close()

	helpers.Check(err)
}
