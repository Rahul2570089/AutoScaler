package models

import "time"

type ScalingActivity struct {
	ID               string
	GroupName        string
	StartTime        time.Time
	EndTime          time.Time
	StatusCode       string
	StatusMessage    string
	Cause            string
	Description      string
	PreviousCapacity int
	DesiredCapacity  int
	ActualCapacity   int
}
