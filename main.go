package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	api "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

const (
	grpcPort = ":8090"
	restPort = ":8080"
)

type server struct {
}

func main() {
	fmt.Println(startMsg(os.Getenv("APP")))

	go listenAndServeGrpc(grpcPort, &server{})
	log.Fatal(listenAndServeRest(restPort, grpcPort))

}

func listenAndServeGrpc(addr string, grpcServer api.ItemServiceServer) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}

	srv := grpc.NewServer()
	api.RegisterItemServiceServer(srv, grpcServer)
	fmt.Printf("gRPC server listening on %v\n", addr)
	return srv.Serve(lis)
}

func listenAndServeRest(addr string, grpcAddr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterItemServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}
	fmt.Printf("REST server listening on %v\n", addr)
	return http.ListenAndServe(addr, mux)
}

func startMsg(app string) string {
	return fmt.Sprintf("Initializing %v server", app)
}

func (s *server) CreateItem(context.Context, *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *server) GetItem(ctx context.Context, in *api.GetItemRequest) (*api.GetItemResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *server) ListItems(context.Context, *api.ListItemsRequest) (*api.ListItemsResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *server) DeleteItem(context.Context, *api.DeleteItemRequest) (*api.DeleteItemResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *server) UpdateItem(context.Context, *api.UpdateItemRequest) (*api.UpdateItemResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}
