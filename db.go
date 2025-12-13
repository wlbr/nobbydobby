package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSink struct {
	config *Config
	//db     *pgx.Conn
	db *pgxpool.Pool
}

func NewPostgresSink(cfg *Config) (*PostgresSink, error) {
	log.Println("Creating new PostgreSQL sink")
	var err error
	s := &PostgresSink{config: cfg}

	//pg connectstring "postgres://user:password@host:port/dbname"
	dbinfo := "postgres://"
	if cfg.PostgreSQL.User != "" {
		dbinfo += cfg.PostgreSQL.User
		if cfg.PostgreSQL.Password != "" {
			dbinfo += ":" + cfg.PostgreSQL.Password
		}
		dbinfo += "@"
	}
	if cfg.PostgreSQL.Host != "" {
		dbinfo += cfg.PostgreSQL.Host
		if cfg.PostgreSQL.Port != "" {
			dbinfo += ":" + cfg.PostgreSQL.Port
		}
	} else {
		log.Println("No PostgresQL host given")
		err = fmt.Errorf("No PostgresQL host given")
	}
	if cfg.PostgreSQL.Database != "" {
		dbinfo += "/" + cfg.PostgreSQL.Database
	} else {
		log.Println("No PostgresQL database name given")
		err = fmt.Errorf("No PostgresQL database name given")
	}
	if err == nil {
		//s.db, err = pgx.Connect(context.Background(), dbinfo)
		s.db, err = pgxpool.New(context.Background(), dbinfo)
		if err != nil {
			log.Println("Cannot open PostgresQL database: %v", err)
		}
		cfg.AddCleanUpFn(func() error {
			log.Println("Cleanup - closing PostgreSQL database connection")
			s.db.Close()
			return nil
		})
		log.Println("Established PostgreSQL database connection")
	} else {
		log.Println("Not creating PostgreSQL sink due to missing connection data")
		return nil, fmt.Errorf("Insufficient PostgreSQL connection data (%s)", err)
	}
	return s, err
}

func (s *PostgresSink) GetUserRegistrations() ([]User, error) {
	log.Println("Getting user registrations from PostgresDB.")
	var users []User

	sql := `
        SELECT id, firstname, lastname, email
        FROM users`

	rows, err := s.db.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *PostgresSink) PutuserRegistration(u *User) error {
	log.Println("Writing user registration to PostgresDB: %+v", u)

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("starting transaction: %w", err)
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(context.Background(), "INSERT INTO users (firstname, lastname,email) VALUES ($1, $3, $2)", u.FirstName, u.LastName, u.Email); err != nil {
		log.Printf("Cannot insert USER into PostgresQL database: %v \n", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
