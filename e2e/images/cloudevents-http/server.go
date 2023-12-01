package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gorilla/mux"
)

const (
	PORT     = 8099
	TLS_PORT = 4333
)

var events []cloudevents.Event

type EmitData struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	events = make([]cloudevents.Event, 0)

	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, receive)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	router := mux.NewRouter()

	router.HandleFunc("/getCloudEvent/{eventreason}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		eventreason := vars["eventreason"]
		filteredEvents := []cloudevents.Event{}

		for _, v := range events {
			emitData := EmitData{}
			if err := json.Unmarshal(v.Data(), &emitData); err != nil {
				fmt.Printf(emitData.Reason)
			} else if emitData.Reason == eventreason {
				filteredEvents = append(filteredEvents, v)
			}
		}

		json.NewEncoder(w).Encode(filteredEvents)
	})

	router.Handle("/", h)

	log.Printf("will listen on :8899\n")
	if err := http.ListenAndServe(":8899", router); err != nil {
		log.Fatalf("unable to start http server, %s", err)
	}
}

func receive(ctx context.Context, event cloudevents.Event) {
	fmt.Printf("Got an Event: %s", event)
	events = append(events, event)
}
