package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	Status         taskdomain.Status         `json:"status"`
	RecurrenceKind taskdomain.RecurrenceKind `json:"recurrence_kind"`
	RecurrenceDays []int                     `json:"recurrence_days"`
}

type taskDTO struct {
	ID             int64                     `json:"id"`
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	Status         taskdomain.Status         `json:"status"`
	RecurrenceKind taskdomain.RecurrenceKind `json:"recurrence_kind"`
	RecurrenceDays []int                     `json:"recurrence_days,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:             task.ID,
		Title:          task.Title,
		Description:    task.Description,
		Status:         task.Status,
		RecurrenceKind: task.RecurrenceKind,
		RecurrenceDays: task.RecurrenceDays,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
	}
}
