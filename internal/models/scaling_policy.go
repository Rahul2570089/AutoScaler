package models

type ScalingPolicy struct {
	ID              string
	GroupName       string
	Direction       PolicyDirection
	AdjustmentType  string
	StepAdjustments []StepAdjustment
	EstimatedWarmup int64
}

type PolicyDirection string

type StepAdjustment struct {
	LowerBound        float64
	UpperBound        float64
	ScalingAdjustment int
}
