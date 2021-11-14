package main

import (
	"encoding/csv"
	"fmt"
	"github.com/OussamaFadlaoui/EnekoEnergie/price_class_selector"
	"github.com/OussamaFadlaoui/EnekoEnergie/types"
	"github.com/OussamaFadlaoui/EnekoEnergie/usage_validator"
	"github.com/OussamaFadlaoui/EnekoEnergie/utils/helpers"
	"io"
	"os"
)

func main() {
	inputFile, err := os.Open(ReadingInputsFilePath)

	helpers.Check(err, "[Error] Could not open inputs file for some reason")

	var reader = csv.NewReader(inputFile)
	var lineCount = 0
	var lastReadingPoint types.ReadingPoint
	var currentReading types.ReadingPoint
	var nextReading types.ReadingPoint
	var calculatedUsageSegment int
	var calculatedUsageSegments = make(map[int][]types.UsageSegment)
	var invalidUsageSegmentIndices = make(map[int][]int)
	var meteringPointCounters = make(map[int]int)
	var flagNextReadingWillBeInvalid bool
	//var results map[int]float64

	for {
		if currentReading == (types.ReadingPoint{}) {
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
		pointId := currentReading.MeteringPointId

		// Flag found for this reading being invalid due to generating an invalid
		// usage segment value
		if flagNextReadingWillBeInvalid {
			lastReadingPoint = nextReading
			flagNextReadingWillBeInvalid = false
			lineCount++
			meteringPointCounters[pointId]++
			continue
		}

		if nextReading.MeteringPointId != currentReading.MeteringPointId {
			lastReadingPoint = nextReading
			meteringPointCounters[pointId]++
			lineCount++
			continue
		}

		calculatedUsageSegment = nextReading.ReadingValue - currentReading.ReadingValue
		pricePerUnit := price_class_selector.GiveCorrectPrice(currentReading.CreatedAt, currentReading.ReadingValue)
		isValidUsageSegment := usage_validator.IsValidUsageSegment(calculatedUsageSegment)

		if isValidUsageSegment {
			calculatedUsageSegments[pointId] = append(calculatedUsageSegments[pointId], types.UsageSegment{
				Usage:        calculatedUsageSegment,
				PricePerUnit: pricePerUnit,
			})
		} else {
			invalidUsageSegmentPlaceholder := types.UsageSegment{
				Usage:        -1,
				PricePerUnit: pricePerUnit,
			}

			if meteringPointCounters[pointId] == 0 {
				// Invalidate the first and second usage segment and mark as invalid. Also skip next reading
				calculatedUsageSegments[pointId] = append(calculatedUsageSegments[pointId], invalidUsageSegmentPlaceholder)
				calculatedUsageSegments[pointId] = append(calculatedUsageSegments[pointId], invalidUsageSegmentPlaceholder)
				invalidUsageSegmentIndices[pointId] = append(invalidUsageSegmentIndices[pointId], 0, 1)
				flagNextReadingWillBeInvalid = true
			} else if meteringPointCounters[pointId] == 2 && len(invalidUsageSegmentIndices[pointId]) == 2 {
				firstSegmentUsage := calculatedUsageSegments[pointId][0].Usage
				linearizationValue := helpers.CapUsageSegment(firstSegmentUsage)
				calculatedUsageSegments[pointId][meteringPointCounters[pointId]] = types.UsageSegment{
					Usage:        linearizationValue,
					PricePerUnit: pricePerUnit,
				}
				calculatedUsageSegments[pointId][meteringPointCounters[pointId]] = types.UsageSegment{
					Usage:        linearizationValue,
					PricePerUnit: pricePerUnit,
				}
			} else if meteringPointCounters[pointId] >= 1 {
				// Invalidate the last usage segment if this wasn't done before
				lastSegmentIndex := meteringPointCounters[pointId] - 1

				// Wasn't invalidated yet
				if helpers.Contains(invalidUsageSegmentIndices[pointId], lastSegmentIndex) {
					invalidUsageSegmentIndices[pointId] = append(invalidUsageSegmentIndices[pointId], lastSegmentIndex)
					calculatedUsageSegments[pointId][lastSegmentIndex].Usage = -1
				}

				//foundValidRefSegments := false
			}
		}

		lastReadingPoint = nextReading
		lineCount++
		meteringPointCounters[pointId]++
	}

	fmt.Printf("%+v\n", calculatedUsageSegments)

	// Patch up all invalid usage segments

	fmt.Println("Finished reading file")
	err = inputFile.Close()

	helpers.Check(err)
}
