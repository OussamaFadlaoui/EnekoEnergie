package types

type ReadingPoint struct {
	MeteringPointId int
	MeteringTypeId int
	ReadingValue int
	CreatedAt    int64
}

type UsageSegment struct {
	Usage        int
	PricePerUnit float64
}
