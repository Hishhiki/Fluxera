package service

import (
	"context"
	"encoding/json"
	"errors"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type ActivityLogService struct {
	logs     *repositories.ActivityLogRepository
	projects *repositories.ProjectRepository
}

func NewActivityLogService(logs *repositories.ActivityLogRepository, projects *repositories.ProjectRepository) *ActivityLogService {
	return &ActivityLogService{logs: logs, projects: projects}
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

	return s.logs.CreateActivityLog(ctx, log)
}

func (s *ActivityLogService) GetByProject(ctx context.Context, projectID, userID int64) ([]*models.ActivityLog, error) {
	_, err := s.projects.GetProjectByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.logs.GetActivityByProjectID(ctx, projectID)
}
