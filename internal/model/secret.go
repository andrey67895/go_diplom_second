package model

import "github.com/google/uuid"

type Secret struct {
	ID       uuid.UUID `db:"id"`
	Encoded  []byte    `db:"encoded"`
	Type     uuid.UUID `db:"type"`
	Metadata string    `db:"metadata"`
}
