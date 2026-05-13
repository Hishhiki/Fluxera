package service

import (
	"context"
	"encoding/json"
	"errors"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type TaskService struct {
	tasks    *repositories.TaskRepository
	projects *repositories.ProjectRepository
	activity ActivityLogger
}
type TaskFilter struct {
	Status string
	Sort   string
}

type ActivityLogger interface {
	Create(ctx context.Context, projectID, userID int64, eventType string, payload json.RawMessage) (*models.ActivityLog, error)
}

func NewTaskService(tasks *repositories.TaskRepository, projects *repositories.ProjectRepository, activity ActivityLogger) *TaskService {
	return &TaskService{
		tasks:    tasks,
		projects: projects,
		activity: activity,
	}
}

func isValidTaskStatus(status string) bool {
	switch status {
	case models.TaskStatusTodo, models.TaskStatusInProgress, models.TaskStatusDone:
		return true
	default:
		return false
	}
}

func isValidTaskPriority(priority string) bool {
	switch priority {
	case models.TaskPriorityNone, models.TaskPriorityLow, models.TaskPriorityMedium, models.TaskPriorityHigh:
		return true
	default:
		return false
	}
}

func isValidTaskSort(sort string) bool {
	switch sort {
	case "created_at_desc", "created_at_asc", "updated_at_desc", "updated_at_asc":
		return true
	default:
		return false
	}
}

func (s *TaskService) getTaskForOwner(ctx context.Context, ownerID, taskID int64) (*models.Task, error) {
	task, err := s.tasks.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	_, err = s.projects.GetProjectByID(ctx, task.ProjectID, ownerID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Create(ctx context.Context, ownerID, projectID int64, title, description, status, priority string) (*models.Task, error) {
	_, err := s.projects.GetProjectByID(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("task title is required")
	}

	description = strings.TrimSpace(description)

	status = strings.TrimSpace(status)
	if status == "" {
		status = models.TaskStatusTodo
	}

	if !isValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}

	priority = strings.TrimSpace(priority)
	if priority == "" {
		priority = models.TaskPriorityNone
	}

	if !isValidTaskPriority(priority) {
		return nil, errors.New("invalid task priority")
	}

	task := &models.Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
	}

	createdTask, err := s.tasks.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(map[string]any{
		"task_id": createdTask.ID,
		"title":   createdTask.Title,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.activity.Create(ctx, projectID, ownerID, "task.created", payload)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (s *TaskService) GetByProject(ctx context.Context, ownerID, projectID int64, filter TaskFilter) ([]*models.Task, error) {
	_, err := s.projects.GetProjectByID(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}

	filter.Status = strings.TrimSpace(filter.Status)
	if filter.Status != "" && !isValidTaskStatus(filter.Status) {
		return nil, errors.New("invalid task status")
	}

	filter.Sort = strings.TrimSpace(filter.Sort)

	if filter.Sort == "" {
		filter.Sort = "created_at_desc"
	}
	if !isValidTaskSort(filter.Sort) {
		return nil, errors.New("invalid task sort")
	}

	return s.tasks.GetTasksByProjectID(ctx, projectID, filter.Status, filter.Sort)
}

func (s *TaskService) UpdateStatus(ctx context.Context, ownerID, taskID int64, status string) (*models.Task, error) {
	existingTask, err := s.getTaskForOwner(ctx, ownerID, taskID)
	if err != nil {
		return nil, err
	}
	status = strings.TrimSpace(status)
	if !isValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}

	updatedTask, err := s.tasks.UpdateTaskStatus(ctx, taskID, status)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(map[string]any{
		"task_id":    updatedTask.ID,
		"old_status": existingTask.Status,
		"new_status": updatedTask.Status,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.activity.Create(ctx, existingTask.ProjectID, ownerID, "task.status_changed", payload)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil

}

func (s *TaskService) Delete(ctx context.Context, ownerID, taskID int64) error {
	task, err := s.getTaskForOwner(ctx, ownerID, taskID)
	if err != nil {
		return err
	}

	return s.tasks.DeleteTask(ctx, taskID, task.ProjectID)
}

func (s *TaskService) Update(ctx context.Context, ownerID, taskID int64, title, description, status, priority string) (*models.Task, error) {
	existingTask, err := s.getTaskForOwner(ctx, ownerID, taskID)
	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	if title == "" {
		title = existingTask.Title
	}

	description = strings.TrimSpace(description)
	if description == "" {
		description = existingTask.Description
	}

	status = strings.TrimSpace(status)
	if status == "" {
		status = existingTask.Status
	}
	if !isValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}

	priority = strings.TrimSpace(priority)
	if priority == "" {
		priority = existingTask.Priority
	}
	if !isValidTaskPriority(priority) {
		return nil, errors.New("invalid task priority")
	}

	task := &models.Task{
		ID:          taskID,
		ProjectID:   existingTask.ProjectID,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
	}

	return s.tasks.UpdateTask(ctx, task)

}
