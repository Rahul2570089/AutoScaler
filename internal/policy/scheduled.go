package policy

import (
	"AutoScaler/internal/models"
	"context"
	"fmt"
	"time"
)

type ScheduledActionEvaluator struct {
	clock Clock
}

func NewScheduledActionEvaluator(clock Clock) *ScheduledActionEvaluator {
	return &ScheduledActionEvaluator{
		clock: clock,
	}
}

func (e *ScheduledActionEvaluator) Evaluate(ctx context.Context, actions []*models.ScheduledAction, group *models.Group) (int, string, bool) {
	now := e.clock.Now()

	for _, action := range actions {
		if action.ScheduledTime.Before(now) || action.ScheduledTime.Equal(now) {
			newCapacity := group.DesiredCapacity

			if action.DesiredCapacity > 0 {
				newCapacity = action.DesiredCapacity
			}
			if action.MinCapacity > 0 && newCapacity < action.MinCapacity {
				newCapacity = action.MinCapacity
			}
			if action.MaxCapacity > 0 && newCapacity > action.MaxCapacity {
				newCapacity = action.MaxCapacity
			}

			reason := fmt.Sprintf("scheduled action %s triggered at %s", action.ID, action.ScheduledTime.Format(time.RFC3339))

			return newCapacity, reason, true

		}
	}
	return group.DesiredCapacity, "no scheduled action due", false
}
