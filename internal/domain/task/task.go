package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceKind string

const (
	RecurrenceNone            RecurrenceKind = "none"
	RecurrenceDaily           RecurrenceKind = "daily"
	RecurrenceWeekly          RecurrenceKind = "weekly"
	RecurrenceMonthlyDates    RecurrenceKind = "monthly_dates"
	RecurrenceMonthlyEvenDays RecurrenceKind = "monthly_even_days"
	RecurrenceMonthlyOddDays  RecurrenceKind = "monthly_odd_days"
)

type Task struct {
	ID             int64          `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Status         Status         `json:"status"`
	RecurrenceKind RecurrenceKind `json:"recurrence_kind"`
	RecurrenceDays []int          `json:"recurrence_days,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (k RecurrenceKind) Valid() bool {
	switch k {
	case RecurrenceNone, RecurrenceDaily, RecurrenceWeekly, RecurrenceMonthlyDates, RecurrenceMonthlyEvenDays, RecurrenceMonthlyOddDays:
		return true
	default:
		return false
	}
}
