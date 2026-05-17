package service

import (
	"context"
	"encoding/json"
	"errors"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
	"time"
)

type ActivityCache interface {
	GetActivityFeed(ctx context.Context, projectID int64) ([]*models.ActivityLog, bool, error)
	SetActivityFeed(ctx context.Context, projectID int64, logs []*models.ActivityLog, ttl time.Duration) error
	DeleteActivityFeed(ctx context.Context, projectID int64) error
}

type ActivityLogService struct {
	logs     *repositories.ActivityLogRepository
	projects *repositories.ProjectRepository
	cache    ActivityCache
}

func NewActivityLogService(logs *repositories.ActivityLogRepository, projects *repositories.ProjectRepository, cache ActivityCache) *ActivityLogService {
	return &ActivityLogService{logs: logs, projects: projects, cache: cache}
}

func (s *ActivityLogService) Create(ctx context.Context, projectID, userID int64, eventType string, payload json.RawMessage) (*models.ActivityLog, error) {
	eventType = strings.TrimSpace(eventType)
	if projectID <= 0 {
		return nil, errors.New("project_id is required")
	}

	if userID <= 0 {
		return nil, errors.New("user_id is required")
	}

	if eventType == "" {
		return nil, errors.New("event type is required")
	}
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}

	if !json.Valid(payload) {
		return nil, errors.New(`invalid payload json`)
	}

	log := &models.ActivityLog{
		ProjectID: projectID,
		UserID:    userID,
		EventType: eventType,
		Payload:   payload,
	}

	createdLog, err := s.logs.CreateActivityLog(ctx, log)
	if err != nil {
		return nil, err
	}

	if err := s.cache.DeleteActivityFeed(ctx, projectID); err != nil {
		return nil, err
	}

	return createdLog, nil
}

func (s *ActivityLogService) GetByProject(ctx context.Context, projectID, userID int64) ([]*models.ActivityLog, error) {
	_, err := s.projects.GetProjectByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	if logs, found, err := s.cache.GetActivityFeed(ctx, projectID); err != nil {
		return nil, err
	} else if found {
		return logs, nil
	}

	logs, err := s.logs.GetActivityByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetActivityFeed(ctx, projectID, logs, 5*time.Minute); err != nil {
		return nil, err
	}

	return logs, nil
}
