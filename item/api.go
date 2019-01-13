package item

import (
	"context"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/model"
	pb "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
)

// Server struct
type Server struct {
	item *Service
}

// NewItemAPI func
func NewItemAPI(service *Service) *Server {
	return &Server{item: service}
}

// CreateItem func
func (s *Server) CreateItem(ctx context.Context, in *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	id, err := s.item.Create(ctx, pbToModel(in.Item))
	if err != nil {
		return nil, err
	}
	return &pb.CreateItemResponse{Id: *id}, nil
}

// GetItem func
func (s *Server) GetItem(ctx context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := s.item.GetByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetItemResponse{
		Item: modelToPb(item),
	}, nil
}

// ListItems func
func (s *Server) ListItems(ctx context.Context, in *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	list, err := s.item.Fetch(ctx, 10)
	if err != nil {
		return nil, err
	}

	var items = []*pb.Item{}
	for _, item := range list {
		items = append(items, modelToPb(item))
	}

	return &pb.ListItemsResponse{Items: items}, nil
}

// UpdateItem func
func (s *Server) UpdateItem(ctx context.Context, in *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	_, err := s.item.Update(ctx, pbToModel(in.Item))
	if err != nil {
		return nil, err
	}

	return &pb.UpdateItemResponse{}, nil
}

// DeleteItem func
func (s *Server) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	ok, err := s.item.Delete(ctx, in.Id)
	if !ok {
		return nil, err
	}
	return &pb.DeleteItemResponse{}, nil
}

func modelToPb(item *model.Item) *pb.Item {
	if item == nil {
		return nil
	}
	return &pb.Item{
		Id:   item.ID,
		Name: item.Name,
	}
}

func pbToModel(item *pb.Item) *model.Item {
	if item == nil {
		return nil
	}
	return &model.Item{
		ID:   item.Id,
		Name: item.Name,
	}
}
