package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func get_api_key() string {
	if key, ok := os.LookupEnv("API_KEY"); ok {
		return key
	} else {
		return "testinghaasapikey"
	}
}

type Handler struct {
	conn dokku.DokkuConn
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	api_key := ""
	api_key2, ok := query["api_key"]
	if !ok {
		// Get from auth header if possible
		if auth_header, ok := r.Header["Authorization"]; ok {
			auth_header2 := auth_header[0]
			if strings.HasPrefix(auth_header2, "Bearer ") {
				api_key = strings.TrimPrefix(auth_header2, "Bearer ")
			}
		}
	} else {
		api_key = api_key2[0]
	}
	if api_key != get_api_key() {
		w.WriteHeader(401)
		return
	}

	args := strings.Split(query.Get("command"), " ")

	res, err := handler.conn.RunCommand(r.Context(), args)
	if err != nil {
		fmt.Fprintf(w, "Error! %s", err)
	} else {
		fmt.Fprintf(w, "Command success:\n%s", res)
	}
}

func get_port() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return "5000"
	}
}

func main() {
	conn, err := dokku.DokkuConnect(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	hand := Handler{conn}
	http.Handle("/", &hand)
	log.Fatal(http.ListenAndServe(":"+get_port(), nil))
}
