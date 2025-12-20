// entities содержит модели сущностей которые хранятся в БД
package entities

// SecureEntity - хранимая в менеджере паролей сущность
type SecureEntity struct {
	ID       string `json:"id"`
	Metadata string `json:"metadata"`
	OwnerID  string `json:"owner_id"`
}
