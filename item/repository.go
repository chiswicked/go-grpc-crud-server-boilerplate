package item

import (
	"context"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/model"
)

// Repository interface
type Repository interface {
	Create(ctx context.Context, item *model.Item) (*string, error)
	GetByID(ctx context.Context, id string) (*model.Item, error)
	Fetch(ctx context.Context, num int64) ([]*model.Item, error)
	Update(ctx context.Context, item *model.Item) (*model.Item, error)
	Delete(ctx context.Context, id string) (bool, error)
}
