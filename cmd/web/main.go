package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/gorilla/sessions"
	"github.com/hunterwilkins2/gordle/internal/hotreload"
	"github.com/hunterwilkins2/gordle/pkg/trie"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Application struct {
	e     *echo.Echo
	words *trie.Trie
}

type config struct {
	port      int
	hotReload bool
}

type Template struct {
	templates *template.Template
	hotReload bool
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	baseData := struct {
		HotReload bool
		Data      interface{}
	}{
		HotReload: t.hotReload,
		Data:      data,
	}
	return t.templates.ExecuteTemplate(w, name, baseData)
}

func main() {
	cfg := config{}
	flag.IntVar(&cfg.port, "port", 4000, "port to run web server")
	flag.BoolVar(&cfg.hotReload, "hot-reload", false, "hot reload browser")
	flag.Parse()

	words, err := trie.NewFromFile("assets/words.txt")
	if err != nil {
		fmt.Printf("error reading word list: %v", err)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"loop": func(from, to int) <-chan int {
			ch := make(chan int)
			go func() {
				for i := from; i < to; i++ {
					ch <- i
				}
				close(ch)
			}()
			return ch
		},
	}
	t := &Template{
		templates: template.Must(template.New("index").Funcs(funcMap).ParseGlob("ui/html/*.html")),
		hotReload: cfg.hotReload,
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = t

	app := &Application{
		e:     e,
		words: words,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	gob.Register([]guess{})
	gob.Register(map[string]bool{})

	e.Logger.SetLevel(log.DEBUG)
	e.Static("/static", "ui/static")
	e.GET("/", app.home)
	e.POST("/guess", app.guess)
	e.GET("/reset", app.reset)

	if cfg.hotReload {
		e.GET("/hot-reload", hotreload.HotReload)
		e.GET("/hot-reload/ready", hotreload.TestAlive)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
