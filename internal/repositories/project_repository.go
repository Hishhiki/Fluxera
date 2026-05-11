package repositories

import (
	"context"
	"database/sql"
	"fluxera/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO projects (owner_id, name, description) VALUES($1, $2, $3) 	RETURNING id, created_at`,
		project.OwnerID,
		project.Name,
		project.Description,
	)
	if err := row.Scan(&project.ID, &project.CreatedAt); err != nil {
		return nil, err
	}

	return project, nil
}
