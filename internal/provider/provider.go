package provider

import (
	"AutoScaler/internal/models"
	"context"
)

type LaunchRequest struct {
	GroupName   string
	LaunchSpec  models.LaunchSpec
	Zone        string
	ClientToken string
}

type Provider interface {
	Launch(ctx context.Context, req LaunchRequest) (models.Instance, error)
	Terminate(ctx context.Context, ids []string) error
	List(ctx context.Context, groupName string) ([]models.Instance, error)
	Zones(ctx context.Context) ([]string, error)
}
