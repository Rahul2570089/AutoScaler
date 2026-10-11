package reconciler

import "AutoScaler/internal/models"

func BalanceAZs(instances []models.Instance, zones []string) map[string]int {
	distribution := make(map[string]int)

	for _, zone := range zones {
		distribution[zone] = 0
	}

	for _, instance := range instances {
		distribution[instance.Zone]++
	}

	return distribution
}
