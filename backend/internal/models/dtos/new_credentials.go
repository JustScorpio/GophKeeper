// dtos содержит объекты для транспортировки данных
package dtos

// NewCredentials - учётные данные (dto - новая запись)
type NewCredentials struct {
	NewSecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}
