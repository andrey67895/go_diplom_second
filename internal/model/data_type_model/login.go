package data_type_model

import "fmt"

type AuthData struct {
	Authentication struct {
		Login    string `json:"Login"`
		Password string `json:"Password"`
	} `json:"Authentication"`
}

func (c *AuthData) toText() string {
	return fmt.Sprintf("%s:%s", c.Authentication.Login, c.Authentication.Password)
}
