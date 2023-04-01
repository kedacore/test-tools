package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS_ADDRESS = nats://nats.<namespace>.svc.cluster.local:4222
	natsAddress := os.Getenv("NATS_ADDRESS")
	if natsAddress == "" {
		log.Fatal("NATS address cannot be empty")
	}
	nc, err := nats.Connect(natsAddress)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %s", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JS context: %s", err)
	}

	switch os.Args[1] {
	case "consumer":
		consumerMessages(js)
	case "publisher":
		publishMessages(js)
	default:
		fmt.Print("invalid mode, use 'consumer' or 'publisher'")
		os.Exit(1)
	}
}

func consumerMessages(js nats.JetStreamContext) {
	sub, err := js.PullSubscribe("ORDERS.*", "PULL_CONSUMER")
	if err != nil {
		log.Fatalf("failed to create pull subscription: %s", err)
	}

	batch := 5
	for i := 0; i < 50; i++ {
		msgs, err := sub.Fetch(batch, nats.MaxWait(2*time.Second))
		if err != nil {
			log.Fatalf("failed to get message: %s", err)
		}
		// Ack messages.
		for _, msg := range msgs {
			msg.AckSync()
		}
		time.Sleep(1 * time.Second)
	}
}

func publishMessages(js nats.JetStreamContext) {
	messagesCount, err := strconv.ParseInt(os.Getenv("NUM_MESSAGES"), 10, 32)
	if err != nil {
		log.Fatalf("number of messages to write should be a number: %s", err.Error())
	}
	for i := 1; i <= int(messagesCount); i++ {
		_, err = js.Publish("ORDERS.scratch", []byte("order "+strconv.Itoa(i)))
		if err == nil {
			log.Printf("published message: %d\n", i)
		}
	}
}
