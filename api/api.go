package api

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/session"
)

const (
	srvAddr = ":8080"
)

type Server struct {
	srv http.Server
	sm  *session.SessionManager
	dm  *database.DatabaseManager
}

func NewServer(dbConnStr, dbName, authConnStr string) *Server {
	s := &Server{
		srv: http.Server{
			Addr: srvAddr,
		},
		sm: session.ConnectSessionManager(authConnStr),
		dm: database.InitDatabaseManager(dbConnStr, dbName),
	}
	s.srv.Handler = s.initApi()
	return s
}

func (s *Server) Close() error {
	_ = s.sm.Close()
	_ = s.dm.Close()

	return s.srv.Shutdown(context.Background())
}

func (s *Server) ListenAndServe() error {
	logger.Infof("starting HTTP server on %s...", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) initApi() http.Handler {
	return middleware.RecoverMiddleware(middleware.AccessLogMiddleware(newApiSubrouter(s.initApiV1())))
}

func newApiSubrouter(hs ...http.Handler) http.Handler {
	api := http.NewServeMux()
	for _, h := range hs {
		api.Handle("/api/", http.StripPrefix("/api", h))
	}
	return api
}
