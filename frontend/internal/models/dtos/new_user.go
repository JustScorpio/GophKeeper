// dtos - объекты для передачи данных
package dtos

// NewUser - пользователь (dto - новая запись)
type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
