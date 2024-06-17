package model

import (
	"database/sql"
	"time"
)

type Habit struct {
	ID            int
	DisplayID     int
	Name          string
	Description   string
	CreatedAt     time.Time
	UserID        int
	LastDoneDate  sql.NullTime
	DaysCompleted int
}

func GetHabitByID(habits []Habit, id int) *Habit {
	for _, habit := range habits {
		if habit.ID == id {
			return &habit
		}
	}
	return nil
}
