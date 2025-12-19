package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

type BoltDatabase struct {
	db *bbolt.DB
}

func NewBoltDatabase(dbPath string) (*BoltDatabase, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &BoltDatabase{db: db}, nil

}

func (b *BoltDatabase) GetUserRegistrations() ([]User, error) {
	var users []User
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user User
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			users = append(users, user)
		}
		return nil
	})
	return users, err
}

func (b *BoltDatabase) PutuserRegistration(u *User) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		id, _ := bucket.NextSequence()
		u.ID = int(id)
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(strconv.Itoa(u.ID)), buf)
	})

	return err

}

func (b *BoltDatabase) Close() {
	b.db.Close()
}
