package main

import (
	"encoding/json"
	"os"
	"sync"
)

// FlatFileDB implements the database interface using a flat file.
type FlatFileDB struct {
	path string
	mu   sync.Mutex
}

// NewFlatFileDB creates a new FlatFileDB.
func NewFlatFileDB(path string) (*FlatFileDB, error) {
	db := &FlatFileDB{
		path: path,
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := db.writeStore([]User{}); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// GetUserRegistrations returns all user registrations from the flat file.
func (db *FlatFileDB) GetUserRegistrations() ([]User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.readStore()
}

// PutuserRegistration adds a new user registration to the flat file.
func (db *FlatFileDB) PutuserRegistration(user *User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	users, err := db.readStore()
	if err != nil {
		return err
	}

	// Simple ID generation
	if len(users) > 0 {
		user.ID = users[len(users)-1].ID + 1
	} else {
		user.ID = 1
	}

	users = append(users, *user)

	return db.writeStore(users)
}

func (db *FlatFileDB) readStore() ([]User, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []User{}, nil
	}
	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *FlatFileDB) writeStore(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, data, 0644)
}

func (db *FlatFileDB) Close() {

}
