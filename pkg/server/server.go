package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/handler"

	"github.com/sirupsen/logrus"
)

// Server decorates a standard echo.Echo server instance
// it may have additional configuration relating to this application instance
type Server struct {
	*echo.Echo
	Version string
}

// NewServer returns a new instance of Server
// it uses a logger middleware using the logger instance provided
func NewServer(echoSrv *echo.Echo, logger *logrus.Entry, version string) Server {
	srv := Server{echoSrv, version}
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	srv.Use(middleware.Recover())
	srv.HideBanner = true
	srv.registerRoutes()
	return srv
}

func (s *Server) registerRoutes() {
	api := s.Group("/api")
	versionHandler := handler.NewVersionHandler(s.Version)
	api.GET("/version", versionHandler.ShowVersion)
}
