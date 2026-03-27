package store

import "database/sql"

type LinkStore struct {
	db *sql.DB
}

type Links struct {
}
