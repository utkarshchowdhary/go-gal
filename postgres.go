package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var globalPostgresDb *sql.DB

func NewPostgresDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func DbInitSchema() error {
	_, err := globalPostgresDb.Exec(
		`
		CREATE TABLE IF NOT EXISTS images (
			id varchar(255) PRIMARY KEY NOT NULL,
			user_id varchar(255) NOT NULL,
			description text NOT NULL,
			location varchar(255) UNIQUE NOT NULL,
			size int NOT NULL,
			created_at timestamp(0) NOT NULL
		 );
		 
		 CREATE INDEX IF NOT EXISTS user_id_idx ON images (user_id);
		`,
	)
	return err
}
