package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fluxera/internal/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(ctx context.Context, addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func activityFeedKey(projectID int64) string {
	return fmt.Sprintf("activity:project:%d", projectID)
}
func userProjectsKey(userID int64) string {
	return fmt.Sprintf("projects:user:%d", userID)
}
func projectTasksKey(projectID int64, status, sort string) string {
	return fmt.Sprintf("tasks:project:%d:status:%s:sort:%s", projectID, status, sort)
}

func (c *RedisCache) GetActivityFeed(ctx context.Context, projectID int64) ([]*models.ActivityLog, bool, error) {
	key := activityFeedKey(projectID)
	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var logs []*models.ActivityLog

	if err := json.Unmarshal([]byte(value), &logs); err != nil {
		return nil, false, err
	}

	return logs, true, nil
}

func (c *RedisCache) SetActivityFeed(ctx context.Context, projectID int64, logs []*models.ActivityLog, ttl time.Duration) error {
	key := activityFeedKey(projectID)
	value, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) DeleteActivityFeed(ctx context.Context, projectID int64) error {
	key := activityFeedKey(projectID)

	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) GetUserProjects(ctx context.Context, userID int64) ([]*models.Project, bool, error) {
	key := userProjectsKey(userID)

	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var projects []*models.Project
	if err := json.Unmarshal([]byte(value), &projects); err != nil {
		return nil, false, err
	}

	return projects, true, nil
}

func (c *RedisCache) SetUserProjects(ctx context.Context, userID int64, projects []*models.Project, ttl time.Duration) error {
	key := userProjectsKey(userID)

	value, err := json.Marshal(projects)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) DeleteUserProjects(ctx context.Context, userID int64) error {
	key := userProjectsKey(userID)

	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) GetProjectTasks(ctx context.Context, projectID int64, status, sort string) ([]*models.Task, bool, error) {
	key := projectTasksKey(projectID, status, sort)

	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var tasks []*models.Task

	if err := json.Unmarshal([]byte(value), &tasks); err != nil {
		return nil, false, err
	}

	return tasks, true, nil
}

func (c *RedisCache) SetProjectTasks(ctx context.Context, projectID int64, status, sort string, tasks []*models.Task, ttl time.Duration) error {
	key := projectTasksKey(projectID, status, sort)

	value, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) DeleteProjectTasks(ctx context.Context, projectID int64) error {
	pattern := fmt.Sprintf("tasks:project:%d:*", projectID)

	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisCache) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		if err := c.client.Expire(ctx, key, ttl).Err(); err != nil {
			return 0, err
		}
	}

	return count, nil
}
