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

	CreateScheduleAction(ctx context.Context, action *models.ScheduledAction) error
	GetScheduleAction(ctx context.Context, id string) (*models.ScheduledAction, error)
	ListScheduleActions(ctx context.Context, name string) ([]*models.ScheduledAction, error)
	DeleteScheduleAction(ctx context.Context, id string) error

	CreateActivity(ctx context.Context, activity *models.ScalingActivity) error
	ListActivities(ctx context.Context, name string) ([]*models.ScalingActivity, error)

	Close() error
}
