package storage

import (
	"AutoScaler/internal/models"
	"context"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type BadgerStorage struct {
	db *badger.DB
}

func NewBadgerStorage(path string) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger database: %w", err)
	}
	return &BadgerStorage{db: db}, nil
}

func (s *BadgerStorage) Close() error {
	return s.db.Close()
}

func (s *BadgerStorage) CreateGroup(ctx context.Context, group *models.Group) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := GroupKey(group.Name)

		_, err := txn.Get(key)
		if err == nil {
			return fmt.Errorf("group already exists: %s", group.Name)
		}
		if err != badger.ErrKeyNotFound {
			return err
		}

		data, err := Marshal(group)
		if err != nil {
			return err
		}

		return txn.Set(key, data)
	})
}

func (s *BadgerStorage) GetGroup(ctx context.Context, name string) (*models.Group, error) {
	var group models.Group
	err := s.db.View(func(txn *badger.Txn) error {
		key := GroupKey(name)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Unmarshal(val, &group)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("group not found: %s", name)
	}

	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *BadgerStorage) UpdateGroup(ctx context.Context, group *models.Group) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := GroupKey(group.Name)

		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return fmt.Errorf("group not found: %s", group.Name)
		}
		if err != nil {
			return err
		}

		var existingGroup models.Group
		err = item.Value(func(val []byte) error {
			return Unmarshal(val, &existingGroup)
		})
		if err != nil {
			return err
		}

		if existingGroup.Version >= group.Version {
			return fmt.Errorf("stale update: current version is %d, received version is %d", existingGroup.Version, group.Version)
		}

		data, err := Marshal(group)
		if err != nil {
			return err
		}

		return txn.Set(key, data)
	})
}

func (s *BadgerStorage) DeleteGroup(ctx context.Context, name string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := GroupKey(name)
		return txn.Delete(key)
	})
}

func (s *BadgerStorage) ListGroups(ctx context.Context) ([]*models.Group, error) {
	var groups []*models.Group
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(GroupPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var group models.Group
				if err := Unmarshal(val, &group); err != nil {
					return err
				}
				groups = append(groups, &group)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return groups, err
}

func (s *BadgerStorage) CreatePolicy(ctx context.Context, id string) (*models.ScalingPolicy, error) {
	var policy models.ScalingPolicy
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(PolicyKey(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Unmarshal(val, &policy)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("policy not found: %s", id)
	}

	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (s *BadgerStorage) GetPolicy(ctx context.Context, id string) (*models.ScalingPolicy, error) {
	var policy models.ScalingPolicy
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(PolicyKey(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Unmarshal(val, &policy)
		})
	})
	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("policy not found: %s", id)
	}
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *BadgerStorage) ListPolicies(ctx context.Context) ([]*models.ScalingPolicy, error) {
	var policies []*models.ScalingPolicy
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PolicyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var policy models.ScalingPolicy
				if err := Unmarshal(val, &policy); err != nil {
					return err
				}
				policies = append(policies, &policy)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return policies, err
}

func (s *BadgerStorage) DeletePolicy(ctx context.Context, id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(PolicyKey(id))
	})
}

func (s *BadgerStorage) CreateScheduledAction(ctx context.Context, action *models.ScheduledAction) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := ScheduledActionKey(action.ID)
		data, err := Marshal(action)
		if err != nil {
			return err
		}
		return txn.Set(key, data)
	})
}

func (s *BadgerStorage) GetScheduledAction(ctx context.Context, id string) (*models.ScheduledAction, error) {
	var action models.ScheduledAction
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(ScheduledActionKey(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Unmarshal(val, &action)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, fmt.Errorf("scheduled action not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (s *BadgerStorage) ListScheduledActions(ctx context.Context, groupName string) ([]*models.ScheduledAction, error) {
	var actions []*models.ScheduledAction
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(ScheduledActionPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var action models.ScheduledAction
				if err := Unmarshal(val, &action); err != nil {
					return err
				}
				if action.GroupName == groupName {
					actions = append(actions, &action)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return actions, err
}

func (s *BadgerStorage) DeleteScheduledAction(ctx context.Context, id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(ScheduledActionKey(id))
	})
}

func (s *BadgerStorage) CreateActivity(ctx context.Context, activity *models.ScalingActivity) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := ActivityKey(activity.ID)
		data, err := Marshal(activity)
		if err != nil {
			return err
		}
		return txn.Set(key, data)
	})
}

func (s *BadgerStorage) ListActivities(ctx context.Context, groupName string, limit int) ([]*models.ScalingActivity, error) {
	var activities []*models.ScalingActivity
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true // Most recent first
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(ActivityPrefix)
		count := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if limit > 0 && count >= limit {
				break
			}
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var activity models.ScalingActivity
				if err := Unmarshal(val, &activity); err != nil {
					return err
				}
				if activity.GroupName == groupName {
					activities = append(activities, &activity)
					count++
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return activities, err
}
