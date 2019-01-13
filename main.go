package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/item"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/errs"
	"github.com/chiswicked/go-grpc-crud-server-boilerplate/service"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Description = "DESCRIPTION Basic CRUD server based on protocol buffers and gRPC in go"
	app.Usage = "USAGE Basic CRUD server based on protocol buffers and gRPC in go"
	app.Version = "0.0.1"
	app.Flags = flags
	app.Action = start

	err := app.Run(os.Args)
	errs.FatalIf("", err)
}

func start(c *cli.Context) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := service.CreateDBConn(c)
	errs.FatalIf("PostgreSQL connection error", err)
	defer db.Close()

	err = service.TestDBConn(
		db,
		c.Int("db-conn-attempts"),
		c.Int("db-conn-interval"),
	)
	errs.PanicIf("PostgreSQL ping error", err)

	itemRepo := item.NewPostgresRepository(db)
	itemService := item.NewService(itemRepo)
	itemAPI := item.NewItemAPI(itemService)

	lsnr := service.StartTCPListener(c.String("service-grpc-port"))

	grpcServer := service.InitGRPCServer(itemAPI)
	go service.StartGRPCServer(grpcServer, lsnr)
	defer grpcServer.GracefulStop()

	gwServer := service.InitGRPCGatewayServer(
		ctx,
		c.String("service-grpc-port"),
		c.String("service-http-port"),
	)
	go service.StartHTTPServer(gwServer)
	defer gwServer.Shutdown(ctx)

	promServer := service.InitPrometheusServer(c.String("prometheus-http-port"))
	go service.StartHTTPServer(promServer)
	defer promServer.Shutdown(ctx)

	waitForShutdown()
}

func waitForShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	log.Printf("Shutting down %s", os.Getenv("APP"))
}
