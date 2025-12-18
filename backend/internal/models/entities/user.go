// entities содержит модели сущностей которые хранятся в БД
package entities

// User - пользователь (В качестве ID выступает Login)
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
