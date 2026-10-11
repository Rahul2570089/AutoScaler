package health

import "AutoScaler/internal/models"

type Checker interface {
	CheckHealth(instance models.Instance) models.HealthStatus
}

type EC2HealthChecker struct{}

func NewEC2HealthChecker() *EC2HealthChecker {
	return &EC2HealthChecker{}
}

func (e *EC2HealthChecker) CheckHealth(instance models.Instance) models.HealthStatus {
	return models.HealthStatus("healthy")
}
