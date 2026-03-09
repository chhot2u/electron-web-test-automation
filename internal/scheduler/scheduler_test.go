package scheduler

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"flowpilot/internal/models"
)

func TestParseCronBasic(t *testing.T) {
	cs, err := ParseCron("* * * * *")
	if err != nil {
		t.Fatalf("ParseCron: %v", err)
	}
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	next := cs.Next(now)
	if next.Minute() != 1 || next.Hour() != 12 {
		t.Errorf("expected 12:01, got %s", next.Format("15:04"))
	}
}

func TestParseCronInterval(t *testing.T) {
	cs, err := ParseCron("*/15 * * * *")
	if err != nil {
		t.Fatalf("ParseCron: %v", err)
	}
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	next := cs.Next(now)
	if next.Minute() != 15 {
		t.Errorf("expected minute 15, got %d", next.Minute())
	}
}

func TestParseCronSpecificTime(t *testing.T) {
	cs, err := ParseCron("30 9 * * 1-5")
	if err != nil {
		t.Fatalf("ParseCron: %v", err)
	}
	now := time.Date(2025, 1, 6, 8, 0, 0, 0, time.UTC)
	next := cs.Next(now)
	if next.Hour() != 9 || next.Minute() != 30 {
		t.Errorf("expected 09:30, got %s", next.Format("15:04"))
	}
}

func TestParseCronInvalid(t *testing.T) {
	_, err := ParseCron("bad")
	if err == nil {
		t.Error("expected error for bad cron")
	}

	_, err = ParseCron("* * *")
	if err == nil {
		t.Error("expected error for 3-field cron")
	}

	_, err = ParseCron("99 * * * *")
	if err == nil {
		t.Error("expected error for out-of-range minute")
	}
}

func TestParseCronRange(t *testing.T) {
	cs, err := ParseCron("0 9-17 * * *")
	if err != nil {
		t.Fatalf("ParseCron: %v", err)
	}
	now := time.Date(2025, 1, 1, 8, 59, 0, 0, time.UTC)
	next := cs.Next(now)
	if next.Hour() != 9 || next.Minute() != 0 {
		t.Errorf("expected 09:00, got %s", next.Format("15:04"))
	}
}

func TestParseCronComma(t *testing.T) {
	cs, err := ParseCron("0,30 * * * *")
	if err != nil {
		t.Fatalf("ParseCron: %v", err)
	}
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	next := cs.Next(now)
	if next.Minute() != 30 {
		t.Errorf("expected minute 30, got %d", next.Minute())
	}
}

type mockScheduleDB struct {
	mu        sync.Mutex
	schedules []models.Schedule
	flows     map[string]*models.RecordedFlow
	updated   map[string]time.Time
}

func (m *mockScheduleDB) ListDueSchedules(ctx context.Context, now time.Time) ([]models.Schedule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var due []models.Schedule
	for _, s := range m.schedules {
		if s.Enabled && s.NextRunAt != nil && !s.NextRunAt.After(now) {
			due = append(due, s)
		}
	}
	return due, nil
}

func (m *mockScheduleDB) UpdateScheduleRun(ctx context.Context, id string, lastRun, nextRun time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updated[id] = lastRun
	return nil
}

func (m *mockScheduleDB) GetRecordedFlow(ctx context.Context, id string) (*models.RecordedFlow, error) {
	if f, ok := m.flows[id]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("flow not found")
}

type mockSubmitter struct {
	mu        sync.Mutex
	submitted []models.Schedule
}

func (m *mockSubmitter) SubmitScheduledTask(ctx context.Context, sched models.Schedule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.submitted = append(m.submitted, sched)
	return nil
}

func TestSchedulerTick(t *testing.T) {
	past := time.Now().Add(-1 * time.Minute)
	db := &mockScheduleDB{
		schedules: []models.Schedule{
			{ID: "s1", Name: "Test", CronExpr: "* * * * *", Enabled: true, NextRunAt: &past},
		},
		flows:   map[string]*models.RecordedFlow{},
		updated: map[string]time.Time{},
	}
	sub := &mockSubmitter{}

	s := New(db, sub, time.Hour)
	s.tick(context.Background())

	sub.mu.Lock()
	defer sub.mu.Unlock()
	if len(sub.submitted) != 1 {
		t.Errorf("expected 1 submitted task, got %d", len(sub.submitted))
	}
}

func TestSchedulerStartStop(t *testing.T) {
	db := &mockScheduleDB{
		schedules: []models.Schedule{},
		flows:     map[string]*models.RecordedFlow{},
		updated:   map[string]time.Time{},
	}
	sub := &mockSubmitter{}

	s := New(db, sub, 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	s.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	s.Stop()
	cancel()
}
