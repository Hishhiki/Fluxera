package repositories

import (
	"context"
	"database/sql"
	"errors"
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

func (r *ProjectRepository) GetProjects(ctx context.Context, ownerID int64) ([]*models.Project, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, owner_id, name, description, created_at FROM projects WHERE owner_id = $1`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*models.Project{}

	for rows.Next() {
		project := &models.Project{}

		err := rows.Scan(
			&project.ID,
			&project.OwnerID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id, ownerID int64) (*models.Project, error) {
	project := &models.Project{}

	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, owner_id, name, description, created_at 
		 FROM projects 
		 WHERE id = $1 AND owner_id = $2`,
		id,
		ownerID,
	)

	err := row.Scan(
		&project.ID,
		&project.OwnerID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id, ownerID int64) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM projects WHERE id = $1 AND owner_id = $2`,
		id,
		ownerID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}
