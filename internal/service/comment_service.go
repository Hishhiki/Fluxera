package service

import (
	"context"
	"errors"
	"fluxera/internal/events"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
	"strings"
)

type CommentService struct {
	comments *repositories.CommentRepository
	tasks    *repositories.TaskRepository
	projects *repositories.ProjectRepository
	events   events.Publisher
}

func NewCommentService(comments *repositories.CommentRepository, tasks *repositories.TaskRepository, projects *repositories.ProjectRepository, events events.Publisher) *CommentService {
	return &CommentService{comments: comments,
		tasks:    tasks,
		projects: projects,
		events:   events}
}

func (s *CommentService) Create(ctx context.Context, taskID, userID int64, content string) (*models.Comment, error) {
	task, err := s.tasks.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	_, err = s.projects.GetProjectByID(ctx, task.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment content is required")
	}
	comment := &models.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}

	createdComment, err := s.comments.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	payload, err := events.NewCommentCreatedPayload(createdComment.ID, createdComment.TaskID)
	if err != nil {
		return nil, err
	}
	err = s.events.Publish(ctx, models.Event{
		ProjectID: task.ProjectID,
		UserID:    userID,
		Type:      models.EventCommentCreated,
		Payload:   payload,
	})
	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

func (s *CommentService) GetByTask(ctx context.Context, taskID, userID int64) ([]*models.Comment, error) {
	task, err := s.tasks.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	_, err = s.projects.GetProjectByID(ctx, task.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	comments, err := s.comments.GetCommentsByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
func (s *CommentService) Delete(ctx context.Context, commentID, userID int64) error {
	comment, err := s.comments.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.UserID == userID {
		return s.comments.DeleteComment(ctx, commentID)
	}

	task, err := s.tasks.GetTaskByID(ctx, comment.TaskID)
	if err != nil {
		return err
	}

	_, err = s.projects.GetProjectByID(ctx, task.ProjectID, userID)
	if err != nil {
		return errors.New("forbidden")
	}

	return s.comments.DeleteComment(ctx, commentID)
}

func (s *CommentService) Update(ctx context.Context, commentID, userID int64, content string) (*models.Comment, error) {
	comment, err := s.comments.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment.UserID != userID {
		return nil, errors.New("forbidden")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment must contain something")
	}
	return s.comments.UpdateComment(ctx, commentID, content)
}
