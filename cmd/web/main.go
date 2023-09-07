package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"

	"github.com/hunterwilkins2/gordle/internal/hotreload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	port := *flag.Int("port", 4000, "port to run web server")
	reload := *flag.Bool("hot-reload", false, "hot reload browser")
	flag.Parse()

	t := &Template{
		templates: template.Must(template.ParseGlob("ui/html/*.html")),
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = t
	e.Static("/static", "ui/static")

	e.Use(middleware.Logger())

	e.GET("/", home)

	if reload {
		e.GET("/hot-reload", hotreload.HotReload)
		e.GET("/hot-reload/ready", hotreload.TestAlive)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
