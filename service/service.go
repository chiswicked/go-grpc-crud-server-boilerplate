package service

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/errs"
	api "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// CreateDBConn func
func CreateDBConn(c *cli.Context) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=60",
		c.String("db-host"),
		c.String("db-port"),
		c.String("db-user"),
		c.String("db-password"),
		c.String("db-name"),
		c.String("db-ssl-mode"),
	)
	return sql.Open("postgres", connStr)
}

// TestDBConn func
func TestDBConn(db *sql.DB, attempts int, interval int) (err error) {
	for i := 0; i < attempts; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	return err
}

// StartTCPListener func
func StartTCPListener(addr string) net.Listener {
	lsnr, err := net.Listen("tcp", addr)
	errs.PanicIf("Server startup failed", err)
	return lsnr
}

// InitGRPCServer func
func InitGRPCServer(srvc api.ItemServiceServer) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_validator.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		)),
	)
	api.RegisterItemServiceServer(grpcServer, srvc)
	grpc_prometheus.Register(grpcServer)
	return grpcServer
}

// StartGRPCServer func
func StartGRPCServer(srvr *grpc.Server, lsnr net.Listener) {
	go func() {
		err := srvr.Serve(lsnr)
		errs.PanicIf("gRPC server error", err)
	}()
}

// InitGRPCGatewayServer func
func InitGRPCGatewayServer(ctx context.Context, grpcAddr string, httpAddr string) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterItemServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	errs.PanicIf("gRPC gateway handler registration error", err)
	return &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
}

// InitPrometheusServer func
func InitPrometheusServer(httpAddr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
}

// StartHTTPServer func
func StartHTTPServer(srvr *http.Server) {
	go func() {
		err := srvr.ListenAndServe()
		if err != http.ErrServerClosed {
			errs.PanicIf("REST server error", err)
		}
	}()
}
