package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
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
