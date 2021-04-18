package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func get_api_key() string {
	if key, ok := os.LookupEnv("API_KEY"); ok {
		return key
	} else {
		return "testinghaasapikey"
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("api_key") != get_api_key() {
		w.WriteHeader(401)
		return
	}
	res, err := dokku.RunCommand(query.Get("command"))
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
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+get_port(), nil))
}
