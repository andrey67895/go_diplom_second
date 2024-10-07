package data_type_model

import "github.com/google/uuid"

type File struct {
	ID       uuid.UUID
	Filename string
	Data     []byte //archive
	Metadata string
}
