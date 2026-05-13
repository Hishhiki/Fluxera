package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fluxera/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func taskOrderBy(sort string) string {
	switch sort {
	case "created_at_asc":
		return "created_at ASC"
	case "updated_at_desc":
		return "updated_at DESC"
	case "updated_at_asc":
		return "updated_at ASC"
	default:
		return "created_at DESC"
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO tasks (project_id, title, description, status, priority)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`,
		task.ProjectID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
	)
	if err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	task := &models.Task{}
	row := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, title, description, status, priority, created_at, updated_at
	FROM tasks
	WHERE id = $1`,
		id,
	)

	err := row.Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) GetTasksByProjectID(ctx context.Context, projectID int64, status, sort string) ([]*models.Task, error) {

	orderBy := taskOrderBy(sort)

	query := `SELECT id, project_id, title, description, status, priority, created_at, updated_at
		FROM tasks
		WHERE project_id = $1`

	args := []any{projectID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY " + orderBy

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*models.Task{}

	for rows.Next() {
		task := &models.Task{}

		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE tasks
		 SET title = $1,
		     description = $2,
		     status = $3,
		     priority = $4,
		     updated_at = NOW()
		 WHERE id = $5
		 RETURNING id, project_id, title, description, status, priority, created_at, updated_at`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ID,
	)

	if err := row.Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id, projectID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 and project_id = $2`, id, projectID)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id int64, status string) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE tasks
		 SET status = $1,
		     updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, project_id, title, description, status, priority, created_at, updated_at`,
		status,
		id,
	)

	task := &models.Task{}

	if err := row.Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return task, nil
}
