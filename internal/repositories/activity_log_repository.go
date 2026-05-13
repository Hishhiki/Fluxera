package repositories

import (
	"context"
	"database/sql"
	"fluxera/internal/models"
)

type ActivityLogRepository struct {
	db *sql.DB
}

func NewActivityLogRepository(db *sql.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) CreateActivityLog(ctx context.Context, log *models.ActivityLog) (*models.ActivityLog, error) {

	row := r.db.QueryRowContext(ctx, `INSERT INTO activity_logs (project_id, user_id, event_type, payload) VALUES($1,$2,$3,$4) RETURNING id, created_at`,
		log.ProjectID,
		log.UserID,
		log.EventType,
		log.Payload,
	)

	if err := row.Scan(&log.ID, &log.CreatedAt); err != nil {
		return nil, err
	}
	return log, nil
}

func (r *ActivityLogRepository) GetActivityByProjectID(ctx context.Context, projectID int64) ([]*models.ActivityLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, project_id, user_id, event_type, payload, created_at FROM activity_logs WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []*models.ActivityLog{}

	for rows.Next() {
		log := &models.ActivityLog{}

		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&log.UserID,
			&log.EventType,
			&log.Payload,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil

}
