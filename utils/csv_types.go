package utils

type ReadingPoint struct {
	MeteringPointId int
	MeteringTypeId int
	ReadingValue int
	CreatedAt    int64
}

type UsageSegmentCollection struct {
	MeteringPointId int
	UsageSegments 	[]int
}
