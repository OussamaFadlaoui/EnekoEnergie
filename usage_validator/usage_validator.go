package usage_validator

func IsValidUsageSegment(usageCalculated int) bool {
	invalidUsageSegmentValidators := []func(int) bool{
		func(num int) bool { return num >= 0 && num <= 100 },
	}

	isValidUsageSegment := true

	for _, validator := range invalidUsageSegmentValidators {
		if !validator(usageCalculated) {
			isValidUsageSegment = false
			break
		}
	}

	return isValidUsageSegment
}
