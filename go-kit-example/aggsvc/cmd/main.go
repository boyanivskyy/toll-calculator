package main

import (
	"net"
	"net/http"
	"os"

	"github.com/boyanivskyy/toll-calculator/go-kit-example/aggsvc/aggendpoint"
	"github.com/boyanivskyy/toll-calculator/go-kit-example/aggsvc/aggservice"
	"github.com/boyanivskyy/toll-calculator/go-kit-example/aggsvc/aggtransport"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	service := aggservice.New(logger)
	endpoints := aggendpoint.New(service, logger)
	httpHandler := aggtransport.NewHTTPHandler(endpoints, logger)

	httpListener, err := net.Listen("tcp", ":3003")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	http.Handle("/metrics", promhttp.Handler())

	logger.Log("transport", "HTTP", "addr", ":3003")
	if err := http.Serve(httpListener, httpHandler); err != nil {
		panic("")
	}
}
