package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"

	"CatPower/api"
)

func main() {
	dbConnStr := flag.String("db_connstr", "postgres@localhost:5432", "postgresql connection string")
	dbName := flag.String("db_name", "postgres", "database name")
	authConnStr := flag.String("auth_connstr", "localhost:8081", "auth-service connection string")
	flag.Parse()

	defer logger.InitLogger().Sync()

	apiService := api.NewServer(*dbConnStr, *dbName, *authConnStr)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs

		logger.Infof("got %s, shutting down server", sig)
		if err := apiService.Close(); err != nil {
			logger.Infof("HTTP server Shutdown err: %s", err)
		}
		close(idleConnsClosed)
	}()

	if err := apiService.ListenAndServe(); err != http.ErrServerClosed {
		logger.Infof("HTTP server ListenAndServe err: %s", err)
	}

	<-idleConnsClosed
}
