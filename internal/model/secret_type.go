package model

import "github.com/google/uuid"

type SecretType struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
