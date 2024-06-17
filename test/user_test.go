package test

import (
	"testing"

	"github.com/ialhejaili/go-habit-tracker/repository"
)

func TestAuthenticateUser(t *testing.T) {
	user, err := repository.AuthenticateUser(TestDB, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatalf("Expected authenticated user, got nil")
	}

	if user.Username != "testuser" {
		t.Fatalf("Expected username %s, got %s", "testuser", user.Username)
	}
}
