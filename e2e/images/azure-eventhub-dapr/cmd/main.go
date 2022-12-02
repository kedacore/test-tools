package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var port string

func init() {
	port = os.Getenv("APP_PORT")
	if port == "" {
		log.Fatalf("missing environment variable %s", "APP_PORT")
	}
}
func main() {
	http.HandleFunc("/azure-eventhub-dapr", func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
		rw.WriteHeader(200)
	})
	http.ListenAndServe(":"+port, nil)
}
