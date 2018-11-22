package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	api "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
	_ "github.com/lib/pq"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

const (
	grpcPort = ":8090"
	restPort = ":8080"

	pgHost     = "localhost"
	pgPort     = "5432"
	pgUsername = "test-username"
	pgPassword = "test-password"
	pgDatabase = "test-database"
	pgSSLmode  = "disable"
)

type server struct {
	db *sql.DB
}

func main() {
	fmt.Println(startMsg(os.Getenv("APP")))
	srv := &server{db: nil}
	var err error
	db, err := createDbConn()
	if err != nil {
		log.Fatalf("PostgreSQL connection error:, %v", err)
	}
	srv.db = db
	fmt.Printf("Connected to PostgreSQL server on %v\n", pgPort)
	go listenAndServeGrpc(grpcPort, srv)
	log.Fatal(listenAndServeRest(restPort, grpcPort))

}

func listenAndServeGrpc(addr string, serviceServer api.ItemServiceServer) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterItemServiceServer(grpcServer, serviceServer)
	fmt.Printf("gRPC server listening on %v\n", addr)
	return grpcServer.Serve(lis)
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

func createDbConn() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable=%s connect_timeout=60",
		pgHost,
		pgPort,
		pgUsername,
		pgPassword,
		pgDatabase,
		pgSSLmode,
	)
	return sql.Open("postgres", connStr)
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
