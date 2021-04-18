package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
