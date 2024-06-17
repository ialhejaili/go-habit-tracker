package repository

import (
	"database/sql"
	"fmt"

	"github.com/ialhejaili/go-habit-tracker/model"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err = db.Exec(query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("error registering user: %v", err)
	}
	return nil
}

func AuthenticateUser(db *sql.DB, username, password string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			//user not found
			return nil, fmt.Errorf("invalid email and/or password")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email and/or password")
	}

	return &user, nil
}

func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM habits WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("error deleting habits: %v", err)
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
