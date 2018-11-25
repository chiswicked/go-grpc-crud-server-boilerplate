package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/errs"
	api "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	grpcPort = ":8090"
	restPort = ":8080"

	pgHost     = "localhost"
	pgPort     = "5432"
	pgUsername = "testusername"
	pgPassword = "testpassword"
	pgDatabase = "testdatabase"
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
	errs.FatalIf("PostgreSQL connection error", err)
	fmt.Printf("Connected to PostgreSQL server on %v\n", pgPort)
	defer db.Close()

	err = db.Ping()
	errs.PanicIf("PostgreSQL ping error", err)
	fmt.Println("PostgreSQL ping ok")
	srv.db = db

	fmt.Printf("Connected to PostgreSQL server on %v\n", pgPort)
	go listenAndServeGrpc(grpcPort, srv)
	go listenAndServeRest(restPort, grpcPort)

	waitForShutdown()
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=60",
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

func (s *server) CreateItem(ctx context.Context, in *api.CreateItemRequest) (*api.CreateItemResponse, error) {
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

	out := &api.CreateItemResponse{Id: uid.String()}
	_, err = s.db.ExecContext(ctx, qry, uid, in.Item.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not insert item into the database: %s", err)
	}

	return out, nil
}

func (s *server) GetItem(ctx context.Context, in *api.GetItemRequest) (*api.GetItemResponse, error) {
	if _, err := uuid.FromString(in.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Not Found")
	}

	qry := `
		SELECT uuid, name
		FROM itemtable
		WHERE uuid = $1;
	`
	out := &api.GetItemResponse{Item: &api.Item{}}
	err := s.db.QueryRowContext(ctx, qry, in.Id).Scan(&out.Item.Id, &out.Item.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "Not Found")
		}
		return nil, grpc.Errorf(codes.Internal, "Could not read item from the database: %s", err)
	}

	return out, nil
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

func waitForShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	log.Printf("Shutting down %s", os.Getenv("APP"))
}
