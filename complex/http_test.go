package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct {
	users []User
}

func (m *mockDB) GetUserRegistrations() ([]User, error) {
	return m.users, nil
}

func (m *mockDB) PutuserRegistration(u *User) error {
	m.users = append(m.users, *u)
	return nil
}

func (m *mockDB) Close() {}

func TestGetRegistrations(t *testing.T) {
	db := &mockDB{
		users: []User{
			{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		},
	}
	cfg := &Config{}
	webserver := &webserver{cfg: cfg, db: db}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webserver.GetRegistrations)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "User: ID=1, FirstName=John, LastName=Doe, Email=john.doe@example.com \n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRegister(t *testing.T) {
	db := &mockDB{}
	cfg := &Config{}
	webserver := &webserver{cfg: cfg, db: db}

	user := &User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webserver.Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Registered!!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if len(db.users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(db.users))
	}

	if db.users[0].FirstName != user.FirstName {
		t.Errorf("Expected FirstName to be %s, got %s", user.FirstName, db.users[0].FirstName)
	}
}
