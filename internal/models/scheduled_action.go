package models

import "time"

type ScheduledAction struct {
	ID              string
	GroupName       string
	ScheduledTime   time.Time
	DesiredCapacity int
	MinCapacity     int
	MaxCapacity     int
}
