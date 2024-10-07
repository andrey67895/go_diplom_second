package data_type_model

import "fmt"

type AuthData struct {
	ID       string
	Login    string
	Password string
	Metadata string
}

func (c *AuthData) toText() string {
	return fmt.Sprintf("%s:%s:%s:%s", c.ID, c.Login, c.Password, c.Metadata)
}
