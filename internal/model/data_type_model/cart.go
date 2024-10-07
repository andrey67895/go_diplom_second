package data_type_model

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/andrey67895/go_diplom_second/internal/helpers"
)

type Cart struct {
	ID       uuid.UUID
	Number   string
	Expired  string
	Holder   string
	CVC      string
	Metadata string
}

func (c *Cart) toText() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", c.ID, c.Number, c.Expired, c.Holder, c.CVC, c.Metadata)
}

func (c *Cart) toJsonEncoded() ([]byte, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return helpers.Compress(marshal), nil
}
