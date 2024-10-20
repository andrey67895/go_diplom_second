package data_type_model

type CardData struct {
	Card struct {
		Number  string `json:"Number"`
		Expired string `json:"Expired"`
		Holder  string `json:"Holder"`
		CVC     string `json:"CVC"`
	} `json:"CreditCardData"`
}
