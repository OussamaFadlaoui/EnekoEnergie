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
	var usageSegments = make(map[int][]types.UsageSegment)
	var invalidUsageSegmentIndices = make(map[int][]int)
	var mtrngPointUsgSegmentCounters = make(map[int]int)
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
		curUsageIndex := mtrngPointUsgSegmentCounters[pointId]

		// Flag found for this reading being invalid due to generating an invalid
		// usage segment value
		if flagNextReadingWillBeInvalid {
			lastReadingPoint = nextReading
			flagNextReadingWillBeInvalid = false
			lineCount++
			mtrngPointUsgSegmentCounters[pointId]++
			continue
		}

		if nextReading.MeteringPointId != currentReading.MeteringPointId {
			lastReadingPoint = nextReading
			mtrngPointUsgSegmentCounters[pointId]++
			lineCount++
			continue
		}

		calculatedUsageSegment = nextReading.ReadingValue - currentReading.ReadingValue
		pricePerUnit := price_class_selector.GiveCorrectPrice(currentReading.CreatedAt, currentReading.ReadingValue)
		isValidUsageSegment := usage_validator.IsValidUsageSegment(calculatedUsageSegment)

		if isValidUsageSegment {
			usageSegments[pointId] = append(usageSegments[pointId], types.UsageSegment{
				Usage:        calculatedUsageSegment,
				PricePerUnit: pricePerUnit,
			})
		} else {
			invalidUsageSegmentPlaceholder := types.UsageSegment{
				Usage:        -1,
				PricePerUnit: pricePerUnit,
			}

			// Current segment is invalid, and it's the first segment calculated
			// Invalidate the first and second usage segment and mark as invalid.
			// Also skip next reading since we assume this would produce another
			// invalid usage segment.
			if curUsageIndex == 0 {
				usageSegments[pointId] = append(usageSegments[pointId], invalidUsageSegmentPlaceholder)
				usageSegments[pointId] = append(usageSegments[pointId], invalidUsageSegmentPlaceholder)
				invalidUsageSegmentIndices[pointId] = append(invalidUsageSegmentIndices[pointId], 0, 1)
				flagNextReadingWillBeInvalid = true
				lastReadingPoint = nextReading
				continue

			} else if curUsageIndex == 2 && len(invalidUsageSegmentIndices[pointId]) == 0 {
				firstSegmentUsage := usageSegments[pointId][0].Usage
				usageSegments[pointId][1] = types.UsageSegment{
					Usage:        helpers.CapUsageSegment(firstSegmentUsage),
					PricePerUnit: pricePerUnit,
				}

				usageSegments[pointId] = append(usageSegments[pointId], types.UsageSegment{
					Usage:        helpers.CapUsageSegment(
						usageSegments[pointId][1].Usage + firstSegmentUsage,
					),
					PricePerUnit: pricePerUnit,
				})

				

				// Current segment is invalid, and we have more than 1 segments calculated before it
			} else if curUsageIndex >= 1 {
				// There are 3 or more usage segments that were calculated before the current
				// one. In this case, we can try to see if the last 2 usage segments before the
				// current one and the invalidated one were valid, so that we can pull a linear
				// line of usage using those values.
				if curUsageIndex >= 3 {
					firstRefSegmentIndex := curUsageIndex - 3
					secondRefSegmentIndex := curUsageIndex - 2

					if helpers.ArrContainsNone(invalidUsageSegmentIndices[pointId], firstRefSegmentIndex, secondRefSegmentIndex) {
						diffBetweenReferenceSegments :=
							usageSegments[pointId][secondRefSegmentIndex].Usage -
								usageSegments[pointId][firstRefSegmentIndex].Usage

						usageSegments[pointId][curUsageIndex-1].Usage = helpers.CapUsageSegment(
							usageSegments[pointId][curUsageIndex-2].Usage + diffBetweenReferenceSegments,
						)

						usageSegments[pointId] = append(usageSegments[pointId], types.UsageSegment{
							Usage: helpers.CapUsageSegment(
								usageSegments[pointId][curUsageIndex-1].Usage + diffBetweenReferenceSegments,
							),
							PricePerUnit: pricePerUnit,
						})
					} else {
						// If either one of the segments before were invalid, we cannot use those
						// values to pull a linear line of usage. So we have to just continue and
						// after marking these as invalid.
						lastSegmentIndex := curUsageIndex - 1
						if !helpers.ArrContains(invalidUsageSegmentIndices[pointId], lastSegmentIndex) {
							invalidUsageSegmentIndices[pointId] = append(
								invalidUsageSegmentIndices[pointId], lastSegmentIndex)
							usageSegments[pointId][lastSegmentIndex].Usage = -1
						}

						usageSegments[pointId] = append(usageSegments[pointId], invalidUsageSegmentPlaceholder)
						invalidUsageSegmentIndices[pointId] = append(invalidUsageSegmentIndices[pointId],
							curUsageIndex)
					}
				}
			}
		}

		lastReadingPoint = nextReading
		lineCount++
		mtrngPointUsgSegmentCounters[pointId]++
	}

	for _, segments := range usageSegments {
		for _, segment := range segments {
			fmt.Printf("%+v\n", segment)
		}
		fmt.Printf("Count: %v\n", len(segments))
	}

	// Patch up all invalid usage segments

	fmt.Println("Finished reading file")
	err = inputFile.Close()

	helpers.Check(err)
}
