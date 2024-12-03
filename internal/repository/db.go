package repository

import (
	"github.com/jmoiron/sqlx"
)

type DB interface {
	sqlx.ExtContext
}
