package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vkdefinition/pkg/authHandlers"
	"vkdefinition/pkg/server"
)

func main() {
	var httpAddr = flag.String("http", ":8085", "http listen address")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"serve", "auth",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	flag.Parse()
	var srv authHandlers.Auth
	{
		srv = authHandlers.NewServiceAuth(logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := server.MakeEndpointsAuth(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := server.NewHTTPServer(context.Background(), endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
