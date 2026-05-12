package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fluxera/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO comments (task_id, user_id, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		comment.TaskID,
		comment.UserID,
		comment.Content,
	)

	if err := row.Scan(&comment.ID, &comment.CreatedAt); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) GetCommentsByTaskID(ctx context.Context, taskID int64) ([]*models.Comment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, task_id, user_id, content, created_at
		 FROM comments
		 WHERE task_id = $1
		 ORDER BY created_at ASC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		comment := &models.Comment{}

		if err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	comment := &models.Comment{}

	row := r.db.QueryRowContext(ctx,
		`SELECT id, task_id, user_id, content, created_at
		 FROM comments
		 WHERE id = $1`,
		id,
	)

	if err := row.Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM comments WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, commentID int64, content string) (*models.Comment, error) {
	comment := &models.Comment{}
	row := r.db.QueryRowContext(ctx,
		`UPDATE comments
	SET content = $1
	WHERE id = $2
	RETURNING id, task_id, user_id, content, created_at`,
		content,
		commentID,
	)

	err := row.Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
