package main

import (
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
	"github.com/spf13/pflag"

	"github.com/hashicorp/go-hclog"
	jaegerGrpc "github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	jaegerCfg "github.com/uber/jaeger-client-go/config"

	"github.com/yandex-cloud/jaeger-ydb-store/plugin"
)

var (
	logger hclog.Logger
)

func init() {
	viper.SetDefault("plugin_http_listen_address", ":15000")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	logger = hclog.New(&hclog.LoggerOptions{
		Name:       "ydb",
		JSONFormat: true,
	})
}

func main() {
	var path = pflag.String("config", "", "full path to configuration file")
	pflag.Parse()

	if len(*path) > 0 {
		var extension = filepath.Ext(*path)
		if len(extension) > 0 {
			extension = extension[1:]
		}
		viper.SetConfigType(extension)

		f, err := os.Open(*path)
		if err != nil {
			log.Fatal("Could not open file", *path)
		}
		err = viper.ReadConfig(f)
		if err != nil {
			log.Fatal("Could not read config file", *path)
		}
	}

	ydbPlugin := plugin.NewYdbStorage()
	ydbPlugin.InitFromViper(viper.GetViper())
	go serveHttp(ydbPlugin.Registry())

	closer := initTracer()
	defer closer.Close()

	logger.Warn("starting plugin")
	jaegerGrpc.Serve(&shared.PluginServices{
		Store:        ydbPlugin,
		ArchiveStore: ydbPlugin,
	})
	logger.Warn("stopped")
}

func serveHttp(gatherer prometheus.Gatherer) {
	mux := http.NewServeMux()
	logger.Warn("serving metrics", "addr", viper.GetString("plugin_http_listen_address"))
	mux.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	if viper.GetBool("ENABLE_PPROF") {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	err := http.ListenAndServe(viper.GetString("plugin_http_listen_address"), mux)
	if err != nil {
		logger.Error("failed to start http listener", "err", err)
		os.Exit(1)
	}
}

func initTracer() io.Closer {
	cfg, err := jaegerCfg.FromEnv()
	if err != nil {
		logger.Error("cfg from env fail", "err", err)
		os.Exit(1)
	}
	closer, err := cfg.InitGlobalTracer("jaeger-ydb-query")
	if err != nil {
		logger.Error("tracer create failed", "err", err)
		os.Exit(1)
	}
	return closer
}
