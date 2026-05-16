package service

import (
	"context"
	"errors"
	"fluxera/internal/events"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type ProjectService struct {
	projects *repositories.ProjectRepository
	events   events.Publisher
}

func NewProjectService(projects *repositories.ProjectRepository, events events.Publisher) *ProjectService {
	return &ProjectService{projects: projects, events: events}
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
	return s.projects.GetProjects(ctx, ownerID)
}

func (s *ProjectService) GetByID(ctx context.Context, id, ownerID int64) (*models.Project, error) {
	return s.projects.GetProjectByID(ctx, id, ownerID)
}

func (s *ProjectService) Delete(ctx context.Context, id, ownerID int64) error {
	return s.projects.DeleteProject(ctx, id, ownerID)
}
