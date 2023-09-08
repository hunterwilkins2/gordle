package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Application) home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
