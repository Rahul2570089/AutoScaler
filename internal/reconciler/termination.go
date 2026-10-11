package reconciler

import (
	"AutoScaler/internal/models"
	"sort"
)

type TerminationPolicy string

const (
	OldesInstance  TerminationPolicy = "OldestInstance"
	NewestInstance TerminationPolicy = "NewestInstance"
	Default        TerminationPolicy = "Default"
)

func SelectInstancesForTermination(instances []models.Instance, policies []string, count int) []models.Instance {
	if len(instances) == 0 || count == 0 {
		return []models.Instance{}
	}

	policy := Default
	if len(policies) > 0 {
		policy = TerminationPolicy(policies[0])
	}

	switch policy {
	case OldesInstance:
		return selectOldest(instances, count)
	case NewestInstance:
		return selectNewest(instances, count)
	default:
		return selectOldest(instances, count)
	}
}

func selectOldest(instances []models.Instance, count int) []models.Instance {
	sorted := make([]models.Instance, len(instances))
	copy(sorted, instances)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].LaunchAt.Before(sorted[j].LaunchAt)
	})

	if count > len(sorted) {
		count = len(sorted)
	}
	return sorted[:count]
}

func selectNewest(instances []models.Instance, count int) []models.Instance {
	sorted := make([]models.Instance, len(instances))
	copy(sorted, instances)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].LaunchAt.After(sorted[j].LaunchAt)
	})

	if count > len(sorted) {
		count = len(sorted)
	}
	return sorted[:count]
}
