package service

import (
	"context"
	"net"
	"net/http"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/errs"
	api "github.com/chiswicked/go-grpc-crud-server-boilerplate/protobuf"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// StartTCPListener func
func StartTCPListener(addr string) net.Listener {
	lsnr, err := net.Listen("tcp", addr)
	errs.PanicIf("Server startup failed", err)
	return lsnr
}

// InitGRPCServer func
func InitGRPCServer(srvc api.ItemServiceServer) *grpc.Server {
	grpcServer := grpc.NewServer()
	api.RegisterItemServiceServer(grpcServer, srvc)
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

// StartGRPCGatewayServer func
func StartGRPCGatewayServer(srvr *http.Server) {
	go func() {
		err := srvr.ListenAndServe()
		if err != http.ErrServerClosed {
			errs.PanicIf("REST server error", err)
		}
	}()
}
