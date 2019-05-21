package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/handler"

	"github.com/sirupsen/logrus"
)

type Server struct {
	*echo.Echo
	Version string
}

func NewServer(echoSrv *echo.Echo, logger *logrus.Entry, version string) Server {
	srv := Server{echoSrv, version}
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	srv.Use(middleware.Recover())
	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	api := s.Group("/api")
	versionHandler := handler.NewVersionHandler(s.Version)
	api.GET("/version", versionHandler.ShowVersion)
}
