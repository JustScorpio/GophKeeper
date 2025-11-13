// repositories contains interface for implementation of such pattern
package repositories

import (
	"context"
)

// Interface implementing pattern "Repository"
type IRepository[T any] interface {
	// GetAll - get all entities
	GetAll(ctx context.Context) ([]T, error)
	// Get - get entity by ID
	Get(ctx context.Context, id string) (*T, error)
	// Create - create entity
	Create(ctx context.Context, IEntity *T) error
	// Update - update entity
	Update(ctx context.Context, IEntity *T) error
	// Delete - delete entity
	Delete(ctx context.Context, id []string, userID string) error

	// CloseConnection - close connection with DB
	CloseConnection()
}
