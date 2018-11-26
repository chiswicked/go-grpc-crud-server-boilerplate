package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/api"
	"github.com/chiswicked/go-grpc-crud-server-boilerplate/errs"
	"github.com/chiswicked/go-grpc-crud-server-boilerplate/service"
	_ "github.com/lib/pq"

	"golang.org/x/net/context"
)

const (
	grpcAddr = ":8090"
	gwAddr   = ":8080"

	pgHost     = "localhost"
	pgPort     = "5432"
	pgUsername = "testusername"
	pgPassword = "testpassword"
	pgDatabase = "testdatabase"
	pgSSLmode  = "disable"
)

func main() {
	fmt.Println(startMsg(os.Getenv("APP")))

	var err error
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	db, err := createDbConn()
	errs.FatalIf("PostgreSQL connection error", err)
	defer db.Close()

	err = db.Ping()
	errs.PanicIf("PostgreSQL ping error", err)

	srv := api.CreateAPI(db)

	lsnr := service.StartTCPListener(grpcAddr)

	grpcServer := service.InitGRPCServer(srv)
	go service.StartGRPCServer(grpcServer, lsnr)
	defer grpcServer.GracefulStop()

	gwServer := service.InitGRPCGatewayServer(ctx, grpcAddr, gwAddr)
	go service.StartGRPCGatewayServer(gwServer)
	defer gwServer.Shutdown(ctx)

	waitForShutdown()
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

func waitForShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	log.Printf("Shutting down %s", os.Getenv("APP"))
}
