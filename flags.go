package main

import (
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "service-grpc-port",
		Usage:  "gRPC port bind address for gRPC",
		EnvVar: "SERVICE_GRPC_PORT",
		Value:  ":8090",
	},
	cli.StringFlag{
		Name:   "service-http-port",
		Usage:  "REST HTTP port",
		EnvVar: "SERVICE_HTTP_PORT",
		Value:  ":8080",
	},
	cli.StringFlag{
		Name:   "prometheus-http-port",
		Usage:  "Prometheus HTTP port to expose metrics",
		EnvVar: "PROMETHEUS_HTTP_PORT",
		Value:  ":8081",
	},
	cli.StringFlag{
		Name:   "db-host",
		Usage:  "Databse host",
		EnvVar: "DB_HOST",
		Value:  "127.0.0.1",
	},
	cli.StringFlag{
		Name:   "db-port",
		Usage:  "Database port",
		EnvVar: "DB_PORT",
		Value:  "5432",
	},
	cli.StringFlag{
		Name:   "db-user",
		Usage:  "Database username",
		EnvVar: "DB_USER",
		Value:  "testusername",
	},
	cli.StringFlag{
		Name:   "db-password",
		Usage:  "Database password",
		EnvVar: "DB_PASSWORD",
		Value:  "testpassword",
	},
	cli.StringFlag{
		Name:   "db-name",
		Usage:  "Database name",
		EnvVar: "DB_NAME",
		Value:  "testdatabase",
	},
	cli.IntFlag{
		Name:   "db-conn-attempts",
		Usage:  "Number of attempts to connect to database before terminating program execution",
		EnvVar: "DB_CONN_ATTEMPTS",
		Value:  10,
	},
	cli.IntFlag{
		Name:   "db-conn-interval",
		Usage:  "Number of milliseconds between attempts to connect to database",
		EnvVar: "DB_CONN_INTERVAL",
		Value:  1000,
	},
	cli.StringFlag{
		Name:   "db-ssl-mode",
		Usage:  "Database name",
		EnvVar: "DB_SSL_MODE",
		Value:  "disable",
	},
}
