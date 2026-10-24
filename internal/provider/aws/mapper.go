package aws

import "AutoScaler/internal/models"

func MapEC2State(state string) models.InstanceState {
	switch state {
	case "pending":
		return models.InstanceState("pending")
	case "running":
		return models.InstanceState("running")
	case "shutting-down", "stopping":
		return models.InstanceState("Shutting-down")
	case "terminated", "stopped":
		return models.InstanceState("terminated")
	default:
		return models.InstanceState("unknowm")

	}
}

func MapInternalState(state models.InstanceState) string {
	switch state {
	case "pending":
		return "pending"
	case "running":
		return "running"
	case "shutting-down":
		return "shutting-down"
	case "terminated":
		return "terminated"
	default:
		return "unknown"

	}
}
