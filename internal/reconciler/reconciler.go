package reconciler

import (
	"AutoScaler/internal/models"
	"AutoScaler/internal/provider"
	"AutoScaler/internal/storage"
	"context"
	"fmt"
	"time"
)

type Reconciler struct {
	provider provider.Provider
	storage  storage.Storage
}

func NewReconciler(provider provider.Provider, storage storage.Storage) *Reconciler {
	return &Reconciler{
		provider: provider,
		storage:  storage,
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, groupName string) error {
	group, err := r.storage.GetGroup(ctx, groupName)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	actualInstances, err := r.provider.List(ctx, groupName)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	currentCapacity := len(actualInstances)
	desiredCapacity := group.DesiredCapacity

	fmt.Printf("Reconcile: group=%s, current=%d, desired=%d\n", groupName, currentCapacity, desiredCapacity)

	if currentCapacity == desiredCapacity {
		return r.logActivity(ctx, groupName, "NoAction",
			"At desired capacity", currentCapacity, desiredCapacity, currentCapacity)
	} else if currentCapacity < desiredCapacity {
		r.scaleUp(ctx, group, actualInstances, desiredCapacity-currentCapacity)
	} else if currentCapacity > desiredCapacity {
		r.scaleDown(ctx, group, actualInstances, currentCapacity-desiredCapacity)
	}

	return nil
}

func (r *Reconciler) scaleUp(ctx context.Context, group *models.Group, actualInstances []models.Instance, count int) error {
	zones, err := r.provider.Zones(ctx)
	if err != nil {
		return fmt.Errorf("failed to get zones: %w", err)
	}

	for i := 0; i < count; i++ {
		zone := zones[i%len(zones)]
		clientToken := fmt.Sprintf("launch=%s-%d", group.Name, time.Now().UnixNano())

		req := provider.LaunchRequest{
			GroupName:   group.Name,
			LaunchSpec:  group.LaunchSpec,
			Zone:        zone,
			ClientToken: clientToken,
		}

		instance, err := r.provider.Launch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to launch instance: %w", err)
		}

		fmt.Printf("Launched instance: %s in zone %s\n", instance.ID, zone)
	}

	return r.logActivity(ctx, group.Name, "ScaleUp", fmt.Sprintf("Launched %d instances", count), len(actualInstances), group.DesiredCapacity, group.DesiredCapacity)
}

func (r *Reconciler) scaleDown(ctx context.Context, group *models.Group, actualInstances []models.Instance, count int) error {
	instancesToTerminate := r.selectInstancesToTerminate(group, actualInstances, count)

	var ids []string
	for _, instance := range instancesToTerminate {
		ids = append(ids, instance.ID)
	}

	err := r.provider.Terminate(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to terminate instances %v: %w", ids, err)
	}
	fmt.Printf("Terminated instances: %v\n", ids)
	return r.logActivity(ctx, group.Name, "ScaleDown", fmt.Sprintf("Terminated %d instances", count), len(actualInstances), group.DesiredCapacity, group.DesiredCapacity)
}

func (r *Reconciler) selectInstancesToTerminate(group *models.Group, actualInstances []models.Instance, count int) []models.Instance {
	return SelectInstancesForTermination(actualInstances, group.TerminationPolicies, count)
}

func (r *Reconciler) logActivity(ctx context.Context, groupName, statusCode, cause string, previous, desired, actual int) error {
	activity := &models.ScalingActivity{
		ID:               fmt.Sprintf("act-%d", time.Now().UnixNano()),
		GroupName:        groupName,
		StartTime:        time.Now(),
		EndTime:          time.Now(),
		StatusCode:       statusCode,
		Cause:            cause,
		Description:      fmt.Sprintf("Capacity changed from %d to %d instances", previous, desired),
		PreviousCapacity: previous,
		DesiredCapacity:  desired,
		ActualCapacity:   actual,
	}
	return r.storage.CreateActivity(ctx, activity)
}
