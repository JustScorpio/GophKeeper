package entities

// Credentials - учётные данные
type Credentials struct {
	SecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}
