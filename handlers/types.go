package handlers

type PlainText struct {
	PlainText string `json:"plain_text"`
}

type SecretID struct {
	Id string `json:"id"`
}

type SecretData struct {
	Data string `json:"data"`
}
