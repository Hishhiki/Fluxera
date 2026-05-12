package service

import (
	"context"
	"errors"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type TaskService struct {
	tasks    *repositories.TaskRepository
	projects *repositories.ProjectRepository
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

func NewTaskService(tasks *repositories.TaskRepository, projects *repositories.ProjectRepository) *TaskService {
	return &TaskService{
		tasks:    tasks,
		projects: projects,
	}
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

	return s.tasks.CreateTask(ctx, task)

}

func (s *TaskService) GetByProject(ctx context.Context, ownerID, projectID int64) ([]*models.Task, error) {
	_, err := s.projects.GetProjectByID(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}

	return s.tasks.GetTasksByProjectID(ctx, projectID)
}

func (s *TaskService) UpdateStatus(ctx context.Context, ownerID, taskID int64, status string) (*models.Task, error) {
	_, err := s.getTaskForOwner(ctx, ownerID, taskID)
	if err != nil {
		return nil, err
	}
	status = strings.TrimSpace(status)
	if !isValidTaskStatus(status) {
		return nil, errors.New("invalid task status")
	}
	return s.tasks.UpdateTaskStatus(ctx, taskID, status)

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
