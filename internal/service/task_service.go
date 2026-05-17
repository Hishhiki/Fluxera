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

type TaskCache interface {
	GetProjectTasks(ctx context.Context, projectID int64, status, sort string) ([]*models.Task, bool, error)
	SetProjectTasks(ctx context.Context, projectID int64, status, sort string, tasks []*models.Task, ttl time.Duration) error
	DeleteProjectTasks(ctx context.Context, projectID int64) error
}

type TaskService struct {
	tasks    *repositories.TaskRepository
	projects *repositories.ProjectRepository
	events   events.Publisher
	cache    TaskCache
}
type TaskFilter struct {
	Status string
	Sort   string
}

func NewTaskService(
	tasks *repositories.TaskRepository,
	projects *repositories.ProjectRepository,
	events events.Publisher,
	cache TaskCache,
) *TaskService {
	return &TaskService{
		tasks:    tasks,
		projects: projects,
		events:   events,
		cache:    cache,
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
	if err := s.cache.DeleteProjectTasks(ctx, projectID); err != nil {
		return nil, err
	}

	payload, err := events.NewTaskCreatedPayload(createdTask.ID, createdTask.Title)
	if err != nil {
		return nil, err
	}

	err = s.events.Publish(ctx, models.Event{
		ProjectID: projectID,
		UserID:    ownerID,
		Type:      models.EventTaskCreated,
		Payload:   payload,
	})

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

	if tasks, found, err := s.cache.GetProjectTasks(ctx, projectID, filter.Status, filter.Sort); err != nil {
		return nil, err
	} else if found {
		return tasks, nil
	}

	tasks, err := s.tasks.GetTasksByProjectID(ctx, projectID, filter.Status, filter.Sort)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetProjectTasks(ctx, projectID, filter.Status, filter.Sort, tasks, 5*time.Minute); err != nil {
		return nil, err
	}

	return tasks, nil
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
	if err := s.cache.DeleteProjectTasks(ctx, existingTask.ProjectID); err != nil {
		return nil, err
	}

	payload, err := events.NewTaskStatusChangedPayload(updatedTask.ID, existingTask.Status, updatedTask.Status)
	if err != nil {
		return nil, err
	}

	err = s.events.Publish(ctx, models.Event{
		ProjectID: existingTask.ProjectID,
		UserID:    ownerID,
		Type:      models.EventTaskStatusChanged,
		Payload:   payload,
	})

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

	if err := s.tasks.DeleteTask(ctx, taskID, task.ProjectID); err != nil {
		return err
	}

	return s.cache.DeleteProjectTasks(ctx, task.ProjectID)
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
	updatedTask, err := s.tasks.UpdateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	if err := s.cache.DeleteProjectTasks(ctx, existingTask.ProjectID); err != nil {
		return nil, err
	}

	payload, err := events.NewTaskUpdatedPayload(
		updatedTask.ID,
		updatedTask.Title,
		updatedTask.Status,
		updatedTask.Priority,
	)
	if err != nil {
		return nil, err
	}

	err = s.events.Publish(ctx, models.Event{
		ProjectID: existingTask.ProjectID,
		UserID:    ownerID,
		Type:      models.EventTaskUpdated,
		Payload:   payload,
	})
	if err != nil {
		return nil, err
	}

	return updatedTask, nil

}
