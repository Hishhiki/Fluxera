package events

import "encoding/json"

func NewTaskCreatedPayload(taskID int64, title string) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"task_id": taskID,
		"title":   title,
	})
	if err != nil {
		return nil, err
	}

	return json.RawMessage(payload), nil
}

func NewTaskStatusChangedPayload(taskID int64, oldStatus, newStatus string) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"task_id":    taskID,
		"old_status": oldStatus,
		"new_status": newStatus,
	})
	if err != nil {
		return nil, err
	}

	return json.RawMessage(payload), nil
}

func NewCommentCreatedPayload(commentID, taskID int64) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"comment_id": commentID,
		"task_id":    taskID,
	})
	if err != nil {
		return nil, err
	}

	return json.RawMessage(payload), nil
}
