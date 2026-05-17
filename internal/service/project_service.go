package service

import (
	"context"
	"errors"
	"fluxera/internal/events"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
	"time"
)

type ProjectCache interface {
	GetUserProjects(ctx context.Context, userID int64) ([]*models.Project, bool, error)
	SetUserProjects(ctx context.Context, userID int64, projects []*models.Project, ttl time.Duration) error
	DeleteUserProjects(ctx context.Context, userID int64) error
}
type ProjectService struct {
	projects *repositories.ProjectRepository
	events   events.Publisher
	cache    ProjectCache
}

func NewProjectService(projects *repositories.ProjectRepository, events events.Publisher, cache ProjectCache) *ProjectService {
	return &ProjectService{projects: projects, events: events, cache: cache}
}

func (s *ProjectService) Create(ctx context.Context, ownerID int64, name, description string) (*models.Project, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, errors.New("project name is required")
	}
	project := &models.Project{
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
	}
	createdProject, err := s.projects.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}
	if err := s.cache.DeleteUserProjects(ctx, ownerID); err != nil {
		return nil, err
	}
	payload, err := events.NewProjectCreatedPayload(
		createdProject.ID,
		createdProject.OwnerID,
		createdProject.Name,
		createdProject.Description,
	)
	if err != nil {
		return nil, err
	}

	err = s.events.Publish(ctx, models.Event{
		ProjectID: createdProject.ID,
		UserID:    createdProject.OwnerID,
		Type:      models.EventProjectCreated,
		Payload:   payload,
	})
	if err != nil {
		return nil, err
	}

	return createdProject, nil

}

func (s *ProjectService) GetAll(ctx context.Context, ownerID int64) ([]*models.Project, error) {
	if projects, found, err := s.cache.GetUserProjects(ctx, ownerID); err != nil {
		return nil, err
	} else if found {
		return projects, nil
	}

	projects, err := s.projects.GetProjects(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetUserProjects(ctx, ownerID, projects, 5*time.Minute); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id, ownerID int64) (*models.Project, error) {
	return s.projects.GetProjectByID(ctx, id, ownerID)
}

func (s *ProjectService) Delete(ctx context.Context, id, ownerID int64) error {
	if err := s.projects.DeleteProject(ctx, id, ownerID); err != nil {
		return err
	}

	return s.cache.DeleteUserProjects(ctx, ownerID)
}
