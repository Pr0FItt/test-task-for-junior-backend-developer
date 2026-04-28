package task

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:          normalized.Title,
		Description:    normalized.Description,
		Status:         normalized.Status,
		RecurrenceKind: normalized.RecurrenceKind,
		RecurrenceDays: normalized.RecurrenceDays,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:             id,
		Title:          normalized.Title,
		Description:    normalized.Description,
		Status:         normalized.Status,
		RecurrenceKind: normalized.RecurrenceKind,
		RecurrenceDays: normalized.RecurrenceDays,
		UpdatedAt:      s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceKind == "" {
		input.RecurrenceKind = taskdomain.RecurrenceNone
	}

	recurrenceDays, err := normalizeRecurrence(input.RecurrenceKind, input.RecurrenceDays)
	if err != nil {
		return CreateInput{}, err
	}
	input.RecurrenceDays = recurrenceDays

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceKind == "" {
		input.RecurrenceKind = taskdomain.RecurrenceNone
	}

	recurrenceDays, err := normalizeRecurrence(input.RecurrenceKind, input.RecurrenceDays)
	if err != nil {
		return UpdateInput{}, err
	}
	input.RecurrenceDays = recurrenceDays

	return input, nil
}

func normalizeRecurrence(kind taskdomain.RecurrenceKind, days []int) ([]int, error) {
	if !kind.Valid() {
		return nil, fmt.Errorf("%w: invalid recurrence kind", ErrInvalidInput)
	}

	normalizedDays := uniqueSortedDays(days)

	switch kind {
	case taskdomain.RecurrenceNone, taskdomain.RecurrenceDaily, taskdomain.RecurrenceMonthlyEvenDays, taskdomain.RecurrenceMonthlyOddDays:
		if len(normalizedDays) > 0 {
			return nil, fmt.Errorf("%w: recurrence_days must be empty for recurrence kind %s", ErrInvalidInput, kind)
		}

		return nil, nil
	case taskdomain.RecurrenceWeekly:
		if len(normalizedDays) == 0 {
			return nil, fmt.Errorf("%w: recurrence_days is required for weekly recurrence", ErrInvalidInput)
		}

		for _, day := range normalizedDays {
			if day < 1 || day > 7 {
				return nil, fmt.Errorf("%w: weekly recurrence_days must be in range 1..7", ErrInvalidInput)
			}
		}

		return normalizedDays, nil
	case taskdomain.RecurrenceMonthlyDates:
		if len(normalizedDays) == 0 {
			return nil, fmt.Errorf("%w: recurrence_days is required for monthly_dates recurrence", ErrInvalidInput)
		}

		for _, day := range normalizedDays {
			if day < 1 || day > 31 {
				return nil, fmt.Errorf("%w: monthly_dates recurrence_days must be in range 1..31", ErrInvalidInput)
			}
		}

		return normalizedDays, nil
	default:
		return nil, fmt.Errorf("%w: invalid recurrence kind", ErrInvalidInput)
	}
}

func uniqueSortedDays(days []int) []int {
	if len(days) == 0 {
		return nil
	}

	uniqueDays := make([]int, 0, len(days))
	for _, day := range days {
		if slices.Contains(uniqueDays, day) {
			continue
		}

		uniqueDays = append(uniqueDays, day)
	}

	slices.Sort(uniqueDays)
	return uniqueDays
}
