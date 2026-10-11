package storage

import (
	"AutoScaler/internal/models"
	"context"
)

type Storage interface {
	CreateGroup(ctx context.Context, group *models.Group) error
	GetGroup(ctx context.Context, name string) (*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, name string) error
	ListGroups(ctx context.Context) ([]*models.Group, error)

	CreatePolicy(ctx context.Context, policy *models.ScalingPolicy) error
	GetPolicy(ctx context.Context, id string) (*models.ScalingPolicy, error)
	ListPolicies(ctx context.Context, name string) ([]*models.ScalingPolicy, error)
	DeletePolicy(ctx context.Context, id string) error

	CreateScheduledAction(ctx context.Context, action *models.ScheduledAction) error
	GetScheduledAction(ctx context.Context, id string) (*models.ScheduledAction, error)
	ListScheduledActions(ctx context.Context, name string) ([]*models.ScheduledAction, error)
	DeleteScheduledAction(ctx context.Context, id string) error

	CreateActivity(ctx context.Context, activity *models.ScalingActivity) error
	ListActivities(ctx context.Context, name string, limit int) ([]*models.ScalingActivity, error)

	Close() error
}
