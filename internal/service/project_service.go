package service

import (
	"context"
	"errors"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type ProjectService struct {
	projects *repositories.ProjectRepository
}

func NewProjectService(projects *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{projects: projects}
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
