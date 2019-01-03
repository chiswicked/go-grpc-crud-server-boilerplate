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
func (s *Server) ListItems(ctx context.Context, in *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	qry := `
		SELECT uuid, name
		FROM itemtable
	`
	var outItems = []*pb.Item{}

	// jozsi := sql.Rows
	rows, err := s.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not read items from the database: %s", err)
	}

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		outItems = append(outItems, &pb.Item{Id: id, Name: name})
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return nil, grpc.Errorf(codes.Internal, "Could not read items from the database: %s", err)
		}
	}

	return &pb.ListItemsResponse{Items: outItems}, nil
}

// DeleteItem func
func (s *Server) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	qry := `
		DELETE FROM itemtable
		WHERE uuid = $1
	`

	res, err := s.db.ExecContext(ctx, qry, in.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not delete item from the database: %s", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Error while deleting item from the database: %s", err)
	}

	// TODO: check rows affected

	return &pb.DeleteItemResponse{}, nil
}

// UpdateItem func
func (s *Server) UpdateItem(context.Context, *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not Implemented")
}
