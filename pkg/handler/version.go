package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type VersionHandler struct {
	version string
}

func NewVersionHandler(version string) VersionHandler {
	return VersionHandler{version}
}

func (v VersionHandler) ShowVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"version": v.version})
}
