// dtos содержит объекты для транспортировки данных
package dtos

// NewSecureEntity - хранимая в менеджере паролей сущность (dto - новая запись)
type NewSecureEntity struct {
	Metadata string `json:"metadata"`
}
