package entities

// SecureEntity - хранимая в менеджере паролей сущность
type SecureEntity struct {
	ID       string `json:"id"`
	Metadata string `json:"metadata"`
}
