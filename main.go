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
	if query.Get("api_key") != get_api_key() {
		w.WriteHeader(401)
		return
	}
	split := strings.Split(query.Get("command"), " ")
	first, last := split[0], split[1:]
	res, err := handler.conn.RunCommand(r.Context(), first, last)
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
