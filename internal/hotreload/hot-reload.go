package hotreload

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func HotReload(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				return
			}

			if err = websocket.Message.Send(ws, "connected"); err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func TestAlive(c echo.Context) error {
	return c.String(http.StatusOK, "hot reload alive")
}
