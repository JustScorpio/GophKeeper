// repositories содержит интерфейс для реализации паттерна "репозиторий"
package repositories

import (
	"context"
)

// Interface реализация паттерна "репозиторий"
type IRepository[Entity any, DTO any] interface {
	// GetAll - получить все сущности
	GetAll(ctx context.Context) ([]Entity, error)
	// Get - получить сущность по ИД
	Get(ctx context.Context, id string) (*Entity, error)
	// Create - создать сущность
	Create(ctx context.Context, dto *DTO) (*Entity, error)
	// Update - изменить сущность
	Update(ctx context.Context, entity *Entity) (*Entity, error)
	// Delete - удалить сущность
	Delete(ctx context.Context, id string) error

	// CloseConnection - закрыть соединение с БД
	CloseConnection()
}
