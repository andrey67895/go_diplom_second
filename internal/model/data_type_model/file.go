package data_type_model

type FileData struct {
	File struct {
		Filename string `json:"Filename"`
		Data     []byte `json:"Any"`
	} `json:"File"`
}
