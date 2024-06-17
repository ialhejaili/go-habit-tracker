package test

import (
	"testing"

	"github.com/ialhejaili/go-habit-tracker/model"
	"github.com/ialhejaili/go-habit-tracker/repository"
)

func TestAddHabit(t *testing.T) {
	habit := &model.Habit{
		Name:        "Test Habit",
		Description: "This is a test habit",
		UserID:      TestUserID,
	}

	err := repository.AddHabit(TestDB, habit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var habitID int
	err = TestDB.QueryRow("SELECT id FROM habits WHERE name = $1", habit.Name).Scan(&habitID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repository.DeleteHabit(TestDB, TestUserID, habit.ID)
	if err != nil {
		t.Fatalf("Error cleaning up test habit: %v", err)
	}
}

func TestListHabits(t *testing.T) {
	habit1 := &model.Habit{
		Name:        "Habit 1",
		Description: "Description 1",
		UserID:      TestUserID,
	}

	habit2 := &model.Habit{
		Name:        "Habit 2",
		Description: "Description 2",
		UserID:      TestUserID,
	}

	err := repository.AddHabit(TestDB, habit1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repository.AddHabit(TestDB, habit2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	habits, err := repository.ListHabits(TestDB, TestUserID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(habits) != 2 {
		t.Fatalf("Expected 2 habits, got %d", len(habits))
	}

	_, err = TestDB.Exec("DELETE FROM habits WHERE user_id = $1", TestUserID)
	if err != nil {
		t.Fatalf("Error cleaning up test habits: %v", err)
	}
}

func TestUpdateHabitDone(t *testing.T) {
	habit := &model.Habit{
		Name:        "Habit 1",
		Description: "Description 1",
		UserID:      TestUserID,
	}

	err := repository.AddHabit(TestDB, habit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repository.UpdateHabitDone(TestDB, habit.ID, habit.DaysCompleted)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var daysCompleted int

	err = TestDB.QueryRow("SELECT days_completed FROM habits WHERE id = $1", habit.ID).Scan(&daysCompleted)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if daysCompleted != 1 {
		t.Fatalf("Expected daysCompleted to be 1, got %v", daysCompleted)
	}

	err = repository.DeleteHabit(TestDB, TestUserID, habit.ID)
	if err != nil {
		t.Fatalf("Error cleaning up test habit: %v", err)
	}
}

func TestListHabitsToDoToday(t *testing.T) {
	habit := &model.Habit{
		Name:        "Habit 1",
		Description: "Description 1",
		UserID:      TestUserID,
	}

	err := repository.AddHabit(TestDB, habit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	habits, err := repository.ListHabitsToDoToday(TestDB, TestUserID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(habits) == 0 {
		t.Fatalf("Expected at least 1 habit, got %d", len(habits))
	}

	_, err = TestDB.Exec("DELETE FROM habits WHERE user_id = $1", TestUserID)

	if err != nil {
		t.Fatalf("Error cleaning up test habits: %v", err)
	}
}

func TestDeleteHabit(t *testing.T) {
	habit := &model.Habit{
		Name:        "Habit 1",
		Description: "Description 1",
		UserID:      TestUserID,
	}

	err := repository.AddHabit(TestDB, habit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = repository.DeleteHabit(TestDB, TestUserID, habit.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int
	err = TestDB.QueryRow("SELECT COUNT(*) FROM habits WHERE id = $1", habit.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected habit to be deleted, but it still exists")
	}
}
