package policy

import (
	"AutoScaler/internal/models"
	"fmt"
)

type StepScaling struct {
	engine *Engine
}

func NewStepScaling(engine *Engine) *StepScaling {
	return &StepScaling{
		engine: engine,
	}
}

func (s *StepScaling) ValidatePolicy(policy *models.ScalingPolicy) error {
	if len(policy.StepAdjustments) == 0 {
		return fmt.Errorf("step adjustments are required")
	}

	lastStep := policy.StepAdjustments[len(policy.StepAdjustments)-1]
	if lastStep.UpperBound != -1 || lastStep.LowerBound != -1 {
		return fmt.Errorf("last step adjustment must be a catch-all (bound = -1)")
	}

	return nil
}
