package models

type Metric struct {
	Name       string
	Namespace  string
	Dimensions map[string]string
	Statistic  string
	Period     int64
}
