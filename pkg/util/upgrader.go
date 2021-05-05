package util

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func MakeWebsocketUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if dev := os.Getenv("HAAS_DEV"); dev != "" {
				// development mode
				return true
			} else {
				// taken from https://github.com/gorilla/websocket/blob/b65e62901fc1c0d968042419e74789f6af455eb9/server.go#L87-L97
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				u, err := url.Parse(origin)
				if err != nil {
					return false
				}

				return u.Host == r.Host
			}
		},
	}
}
