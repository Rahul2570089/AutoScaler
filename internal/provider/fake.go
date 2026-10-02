package provider

import (
	"AutoScaler/internal/models"
	"context"
	"fmt"
	"sync"
	"time"
)

type FakeProvider struct {
	mu        sync.RWMutex
	instances map[string]models.Instance
	nextID    int
	zones     []string
}

func NewFakeProvider(zones []string) *FakeProvider {
	return &FakeProvider{
		instances: make(map[string]models.Instance),
		zones:     zones,
	}
}

func (f *FakeProvider) Launch(ctx context.Context, req LaunchRequest) (models.Instance, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.nextID++
	id := fmt.Sprintf("i-%d", f.nextID)

	instance := models.Instance{
		ID:        id,
		GroupName: req.GroupName,
		Zone:      req.Zone,
		State:     models.InstanceState("pending"),
		Health:    models.HealthStatus("unknown"),
		LaunchAt:  time.Now(),
	}

	f.instances[id] = instance
	return instance, nil
}

func (f *FakeProvider) Terminate(ctx context.Context, ids []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, id := range ids {
		if inst, exists := f.instances[id]; exists {
			inst.State = models.InstanceState("terminating")
			f.instances[id] = inst
		}
	}
	return nil
}

func (f *FakeProvider) List(ctx context.Context, groupName string) ([]models.Instance, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var instances []models.Instance
	for _, instance := range f.instances {
		if instance.GroupName == groupName {
			instances = append(instances, instance)
		}
	}
	return instances, nil
}

func (f *FakeProvider) Zones(ctx context.Context) ([]string, error) {
	return f.zones, nil
}
