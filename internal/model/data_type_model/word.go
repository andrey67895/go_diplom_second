package data_type_model

import (
	"fmt"

	"github.com/google/uuid"
)

type Word struct {
	ID       uuid.UUID
	Text     string
	Metadata string
}

func (c *Word) toText() string {
	return fmt.Sprintf("%s:%s:%s", c.ID, c.Text, c.Metadata)
}
