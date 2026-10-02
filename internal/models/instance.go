package models

import "time"

type Instance struct {
	ID        string
	GroupName string
	Zone      string
	State     InstanceState
	Health    HealthStatus
	LaunchAt  time.Time
}

type InstanceState string
type HealthStatus string
