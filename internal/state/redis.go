package state

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(addr string) *Store {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Store{client: client}
}

func (s *Store) RecordAgentActivity(ctx context.Context, swarmID, agentID string) error {
	agentsKey := fmt.Sprintf("swarm:%s:agents", swarmID)
	if err := s.client.SAdd(ctx, agentsKey, agentID).Err(); err != nil {
		return fmt.Errorf("state: failed to add agent to swarm set: %w", err)
	}
	agentKey := fmt.Sprintf("swarm:%s:agent:%s", swarmID, agentID)
	if err := s.client.HSet(ctx, agentKey, "last_seen", time.Now().UTC().Format(time.RFC3339)).Err(); err != nil {
		return fmt.Errorf("state: failed to set agent last_seen: %w", err)
	}
	if err := s.client.HIncrBy(ctx, agentKey, "event_count", 1).Err(); err != nil {
		return fmt.Errorf("state: failed to increment event_count: %w", err)
	}

	return nil
}

func (s *Store) PushToWindow(ctx context.Context, swarmID, entry string, maxSize int64) error {
	windowKey := fmt.Sprintf("swarm:%s:window", swarmID)
	if err := s.client.LPush(ctx, windowKey, entry).Err(); err != nil {
		return fmt.Errorf("state: failed to push to window: %w", err)
	}
	if err := s.client.LTrim(ctx, windowKey, 0, maxSize-1).Err(); err != nil {
		return fmt.Errorf("state: failed to trim window: %w", err)
	}
	return nil
}

func (s *Store) GetWindow(ctx context.Context, swarmID string) ([]string, error) {
	windowKey := fmt.Sprintf("swarm:%s:window", swarmID)
	entries, err := s.client.LRange(ctx, windowKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("state: failed to read window: %w", err)
	}
	return entries, nil
}

func (s *Store) GetCostEWMA(ctx context.Context, swarmID string) (float64, error) {
	key := fmt.Sprintf("swarm:%s:cost_ewma", swarmID)
	val, err := s.client.Get(ctx, key).Float64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("state: failed to get cost ewma: %w", err)
	}
	return val, nil
}

func (s *Store) SetCostEWMA(ctx context.Context, swarmID string, ewma float64) error {
	key := fmt.Sprintf("swarm:%s:cost_ewma", swarmID)
	if err := s.client.Set(ctx, key, ewma, 0).Err(); err != nil {
		return fmt.Errorf("state: failed to set cost ewma: %w", err)
	}

	return nil
}

func (s *Store) RecordCompletion(ctx context.Context, swarmID, agentID string) error {
	key := fmt.Sprintf("swarm:%s:completions", swarmID)
	if err := s.client.RPush(ctx, key, agentID).Err(); err != nil {
		return fmt.Errorf("state: failed to record completion: %w", err)
	}
	return nil
}

func (s *Store) GetCompletions(ctx context.Context, swarmID string) ([]string, error) {
	key := fmt.Sprintf("swarm:%s:completions", swarmID)
	entries, err := s.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("state: failed to get completions: %w", err)
	}
	return entries, nil
}
