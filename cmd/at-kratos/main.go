package main

import (
	"encoding/json"
	"flag"
	"github.com/go-kratos/kratos/v2/config/file"
	"os"

	"at-kratos/internal/conf"
	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	apolloC := config.New(
		config.WithSource(
			apollo.NewSource(
				apollo.WithEndpoint(bc.Dependence.Apollo.Endpoint),
				apollo.WithCluster(bc.Dependence.Apollo.Cluster),
				apollo.WithNamespace(bc.Dependence.Apollo.Namespace),
				apollo.WithAppID(bc.Dependence.Apollo.AppId),
			),
		),
	)
	if err := apolloC.Load(); err != nil {
		panic(err)
	}
	var dataC conf.Data
	var tempC string
	tempC, _ = apolloC.Value("application.data").String()
	if err := json.Unmarshal([]byte(tempC), &dataC); err != nil {
		panic(err)
	}
	app, cleanup, err := wireApp(bc.Server, &dataC, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
