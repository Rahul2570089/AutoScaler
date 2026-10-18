package policy

import (
	"AutoScaler/internal/models"
	"context"
	"fmt"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (r *RealClock) Now() time.Time {
	return time.Now()
}

type Engine struct {
	clock Clock
}

func NewEngine(clock Clock) *Engine {
	return &Engine{
		clock: clock,
	}
}

func (e *Engine) Evaluate(ctx context.Context, group *models.Group,
	policies []*models.ScalingPolicy, metricValue float64,
	alarmThreshold float64, currentCapacity int) (int, string, error) {

	if len(policies) == 0 {
		return currentCapacity, "no policies attached", nil
	}

	var scaleOutPolicies []*models.ScalingPolicy
	var scaleInPolicies []*models.ScalingPolicy

	for _, policy := range policies {
		if policy.Direction == models.PolicyDirection("ScaleOut") {
			scaleOutPolicies = append(scaleOutPolicies, policy)
		} else if policy.Direction == models.PolicyDirection("ScaleIn") {
			scaleInPolicies = append(scaleInPolicies, policy)
		}
	}

	metricExceedsThreshold := metricValue > alarmThreshold
	var desiredCapacity int
	var reason string

	if metricExceedsThreshold {
		desiredCapacity, reason = e.evaluateScaleOut(scaleOutPolicies, metricValue, alarmThreshold, currentCapacity, group)
	} else {
		desiredCapacity, reason = e.evaluateScaleIn(scaleInPolicies, metricValue, alarmThreshold, currentCapacity, group)
	}

	if desiredCapacity < group.MinCapacity {
		desiredCapacity = group.MinCapacity
		reason = fmt.Sprintf("%s (clamped to min capacity %d)", reason, group.MinCapacity)
	}

	if desiredCapacity > group.MaxCapacity {
		desiredCapacity = group.MaxCapacity
		reason = fmt.Sprintf("%s (clamped to max capacity %d)", reason, group.MaxCapacity)
	}

	return desiredCapacity, reason, nil
}

func (e *Engine) evaluateScaleOut(policies []*models.ScalingPolicy, metricValue float64,
	alarmThreshold float64, currentCapacity int, group *models.Group) (int, string) {

	if len(policies) == 0 {
		return currentCapacity, "no scale out policies"
	}

	maxCapacity := currentCapacity
	maxReason := "no scale out policies"

	for _, policy := range policies {
		capacity, reason := e.evluateStepAdjustments(policy, metricValue, alarmThreshold, currentCapacity, group)
		if capacity > maxCapacity {
			maxCapacity = capacity
			maxReason = reason
		}
	}

	return maxCapacity, maxReason
}

func (e *Engine) evaluateScaleIn(policies []*models.ScalingPolicy, metricValue float64,
	alarmThreshold float64, currentCapacity int, group *models.Group) (int, string) {

	if len(policies) == 0 {
		return currentCapacity, "no scale in policies"
	}

	firstCapacity := -1
	firstReason := ""
	allAgree := true

	for _, policy := range policies {
		capacity, reason := e.evluateStepAdjustments(policy, metricValue, alarmThreshold, currentCapacity, group)

		if firstCapacity == -1 {
			firstCapacity = capacity
			firstReason = reason
		} else if capacity != firstCapacity {
			allAgree = false
			break
		}
	}

	if allAgree && firstCapacity != -1 {
		return firstCapacity, firstReason
	}

	return currentCapacity, "scale-in policies do not agree"
}

func (e *Engine) evluateStepAdjustments(policy *models.ScalingPolicy, metricValue float64,
	alarmThreshold float64, currentCapacity int, group *models.Group) (int, string) {

	metricOffset := metricValue - alarmThreshold

	for _, step := range policy.StepAdjustments {
		if e.matchesStep(metricOffset, step) {
			adjustment := step.ScalingAdjustment

			var newCapacity int
			var reason string

			switch policy.AdjustmentType {
			case "ChangeInCapacity":
				newCapacity = currentCapacity + adjustment
				reason = fmt.Sprintf("step adjustment: %s%d (metric: %.2f, threshold: %.2f, offset: %.2f)", sign(adjustment), adjustment, metricValue, alarmThreshold, metricOffset)
			case "ExactCapacity":
				newCapacity = adjustment
				reason = fmt.Sprintf("exact capacity: %d", adjustment)
			default:
				newCapacity = currentCapacity
				reason = fmt.Sprintf("unknown adjustment type: %s", policy.AdjustmentType)
			}

			return newCapacity, reason
		}
	}

	return currentCapacity, "no matching step adjustment"
}

func (e *Engine) matchesStep(metricOffset float64, step models.StepAdjustment) bool {
	if step.LowerBound == -1 && step.UpperBound == -1 {
		return true
	}

	if step.LowerBound != -1 && metricOffset < step.LowerBound {
		return false
	}

	if step.UpperBound != -1 && metricOffset >= step.UpperBound {
		return false
	}

	return true
}

func sign(n int) string {
	if n < 0 {
		return ""
	}
	return "+"
}
