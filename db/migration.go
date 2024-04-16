package db

import "github.com/jmoiron/sqlx"

type Migration struct {
	ID   string
	Up   func(*sqlx.Tx) error
	Down func(*sqlx.Tx) error
}
