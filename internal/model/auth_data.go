package model

import "github.com/google/uuid"

type AuthData struct {
	ID             uuid.UUID `db:"id"`
	Login          string    `db:"login"`
	HashPass       string    `db:"hash_pass"`
	HashPassmaster string    `db:"hash_pass_master"`
}
