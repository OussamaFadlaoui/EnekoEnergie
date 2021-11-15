package helpers

func CapUsageSegment(calculatedUsage int) int {
	if calculatedUsage > 100 {
		return 100
	} else if calculatedUsage < 0 {
		return 0
	}
	return calculatedUsage
}
