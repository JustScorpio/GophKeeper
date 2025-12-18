// entities содержит модели сущностей которые хранятся в БД
package entities

type Hashable interface {
	GetHash() string
}

func Equals(a, b Hashable) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.GetHash() == b.GetHash()
}
