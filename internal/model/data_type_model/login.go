package data_type_model

type AuthData struct {
	Authentication struct {
		Login    string `json:"Login"`
		Password string `json:"Password"`
	} `json:"Authentication"`
}
