package main

type database interface {
	GetUserRegistrations() ([]User, error)
	PutuserRegistration(*User) error
}
