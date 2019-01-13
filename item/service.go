package item

import (
	"context"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/model"
)

// Service struct
type Service struct {
	repo Repository
}

// NewService func
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// Create func
func (s *Service) Create(ctx context.Context, item *model.Item) (*string, error) {
	id, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetByID func
func (s *Service) GetByID(ctx context.Context, id string) (*model.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Fetch func
func (s *Service) Fetch(ctx context.Context, num int64) ([]*model.Item, error) {
	list, err := s.repo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Update func
func (s *Service) Update(ctx context.Context, item *model.Item) (*model.Item, error) {
	item, err := s.repo.Update(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

//Delete func
func (s *Service) Delete(ctx context.Context, id string) (bool, error) {
	_, err := s.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
