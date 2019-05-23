package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// VersionHandler is used to expose the version to the version endpoint
type VersionHandler struct {
	version string
}

// NewVersionHandler returns an instance of VersionHandler
// it initialises it with the version provided
func NewVersionHandler(version string) VersionHandler {
	return VersionHandler{version}
}

// ShowVersion is the handler for GET /version
func (v VersionHandler) ShowVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"version": v.version})
}
