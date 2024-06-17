package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ialhejaili/go-habit-tracker/model"
)

func AddHabit(db *sql.DB, habit *model.Habit) error {
	query := `
		INSERT INTO habits (name, description, user_id, last_done_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	err := db.QueryRow(query, habit.Name, habit.Description, habit.UserID, habit.LastDoneDate).Scan(&habit.ID)
	if err != nil {
		return fmt.Errorf("error adding habit: %v", err)
	}
	return nil
}

func ListHabits(db *sql.DB, userID int) ([]model.Habit, error) {
	query := `SELECT id, name, description, created_at, last_done_date, days_completed 
	FROM habits 
	WHERE user_id = $1 
	ORDER BY id`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error listing habits: %v", err)
	}
	defer rows.Close()

	habits := make([]model.Habit, 0)
	displayID := 1

	for rows.Next() {
		var habit model.Habit
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.CreatedAt, &habit.LastDoneDate, &habit.DaysCompleted)
		if err != nil {
			return nil, fmt.Errorf("error scanning habit row: %v", err)
		}

		habit.DisplayID = displayID
		displayID++

		habits = append(habits, habit)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return habits, nil
}

func UpdateHabitDone(db *sql.DB, habitID int, daysCompleted int) error {
	today := time.Now().Format("2006-01-02")
	daysCompleted++
	query := `UPDATE habits SET last_done_date = $1, days_completed = $2 WHERE id = $3`
	_, err := db.Exec(query, today, daysCompleted, habitID)
	if err != nil {
		return fmt.Errorf("error updating habit done status: %v", err)
	}
	return nil
}

func ListHabitsToDoToday(db *sql.DB, userID int) ([]model.Habit, error) {
	today := time.Now().Format("2006-01-02")
	query := `SELECT id, name, description, created_at, last_done_date, days_completed 
	FROM habits 
	WHERE user_id = $1 AND (last_done_date IS NULL OR last_done_date < $2) 
	ORDER BY id`
	rows, err := db.Query(query, userID, today)
	if err != nil {
		return nil, fmt.Errorf("error listing habits to do today: %v", err)
	}
	defer rows.Close()

	habits := make([]model.Habit, 0)

	for rows.Next() {
		var habit model.Habit
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.CreatedAt, &habit.LastDoneDate, &habit.DaysCompleted)
		if err != nil {
			return nil, fmt.Errorf("error scanning habit row: %v", err)
		}

		habits = append(habits, habit)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return habits, nil
}

func DeleteHabit(db *sql.DB, userID int, habitID int) error {
	query := `DELETE FROM habits WHERE id = $1 AND user_id = $2`
	result, err := db.Exec(query, habitID, userID)
	if err != nil {
		return fmt.Errorf("error deleting habit: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}
