package api

import (
	"context"
	"database/sql"

	pb "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server struct
type Server struct {
	db *sql.DB
}

// CreateAPI func
func CreateAPI(db *sql.DB) *Server {
	return &Server{db: db}
}

// CreateItem func
func (s *Server) CreateItem(ctx context.Context, in *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	if len(in.Item.Name) <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid Argument")
	}

	qry := `
		INSERT INTO itemtable (uuid, name)
		VALUES ($1, $2);
	`
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not insert item into the database: %s", err)
	}

	out := &pb.CreateItemResponse{Id: uid.String()}
	_, err = s.db.ExecContext(ctx, qry, uid, in.Item.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not insert item into the database: %s", err)
	}

	return out, nil
}

// GetItem func
func (s *Server) GetItem(ctx context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	if _, err := uuid.FromString(in.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Not Found")
	}

	qry := `
		SELECT uuid, name
		FROM itemtable
		WHERE uuid = $1;
	`
	out := &pb.GetItemResponse{Item: &pb.Item{}}
	err := s.db.QueryRowContext(ctx, qry, in.Id).Scan(&out.Item.Id, &out.Item.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "Not Found")
		}
		return nil, grpc.Errorf(codes.Internal, "Could not read item from the database: %s", err)
	}

	return out, nil
}

// ListItems func
func (s *Server) ListItems(context.Context, *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not Implemented")
}

// DeleteItem func
func (s *Server) DeleteItem(context.Context, *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not Implemented")
}

// UpdateItem func
func (s *Server) UpdateItem(context.Context, *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not Implemented")
}