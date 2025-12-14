package main

import (
	"context"
	"testing"
)

func TestGetUserRegistrations(t *testing.T) {
	cfg := &Config{}
	cfg.PostgreSQL.Host = "localhost"
	cfg.PostgreSQL.Port = "5432"
	cfg.PostgreSQL.Database = "felix"
	cfg.PostgreSQL.User = "felixapp"
	cfg.PostgreSQL.Password = ""

	db, err := NewPostgresSink(cfg)
	if err != nil {
		t.Fatalf("Could not get db connection: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.db.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatalf("Could not clean up users table: %v", err)
		}
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	err = db.PutuserRegistration(user)
	if err != nil {
		t.Fatalf("Could not put user registration: %v", err)
	}

	users, err := db.GetUserRegistrations()
	if err != nil {
		t.Fatalf("Could not get user registrations: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}

	if users[0].FirstName != user.FirstName {
		t.Errorf("Expected FirstName to be %s, got %s", user.FirstName, users[0].FirstName)
	}

	if users[0].LastName != user.LastName {
		t.Errorf("Expected LastName to be %s, got %s", user.LastName, users[0].LastName)
	}

	if users[0].Email != user.Email {
		t.Errorf("Expected Email to be %s, got %s", user.Email, users[0].Email)
	}
}
