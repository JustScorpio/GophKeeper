package models

// Credentials - credentials to store
type Credentials struct {
	ID       string `json:"id"`
	Metadata string `json:"metadata"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
