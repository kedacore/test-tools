package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-stomp/stomp/v3"
)

const (
	destination = "test"
)

func main() {
	actionType := os.Args[1]
	if actionType == "producer" {
		publishMessages()
	}
	if actionType == "consumer" {
		consumeMessages()
	}

}

func publishMessages() {
	messages := 100
	if val, ok := os.LookupEnv("ARTEMIS_MESSAGE_COUNT"); ok {
		m, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		messages = m

	}
	log.Printf("Publishing %v messages", messages)
	conn, err := getConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Disconnect()
	for i := 0; i < messages; i++ {
		log.Printf("Publishing %v of %v", (i + 1), messages)
		msg := fmt.Sprintf("Message %v", i)
		m, err := encodeGob(msg)
		if err != nil {
			panic(fmt.Errorf("failed to encode message: %v: %v", msg, err))
		}
		err = conn.Send(destination, "text/plain", m, stomp.SendOpt.Header("destination-type", "ANYCAST"))
		if err != nil {
			panic(fmt.Errorf("could not send to destination %s: %v", destination, err))
		}
	}
}

func consumeMessages() {
	log.Print("Reading messages")
	sleep := 100
	if val, ok := os.LookupEnv("ARTEMIS_MESSAGE_SLEEP_MS"); ok {
		s, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		sleep = s
	}

	conn, err := getConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Disconnect()
	sub, err := conn.Subscribe(destination, stomp.AckAuto, stomp.SubscribeOpt.Header("subscription-type", "ANYCAST"))
	if err != nil {
		panic(fmt.Errorf("could not subscribe to queue %s: %v", destination, err))
	}
	for {
		msg := <-sub.C
		if msg.Err != nil {
			panic(fmt.Errorf("received an error: %v", msg.Err))
		}
		m, err := decodeGob[string](msg.Body)
		if err != nil {
			panic(fmt.Errorf("failed to decode message: %v: %v", msg.Header, err))
		}
		fmt.Println(*m)
		time.Sleep(time.Duration(sleep * int(time.Millisecond)))
	}
}

func getConnection() (*stomp.Conn, error) {
	return stomp.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("ARTEMIS_SERVER_HOST"), os.Getenv("ARTEMIS_SERVER_PORT")), stomp.ConnOpt.Login(os.Getenv("ARTEMIS_USERNAME"), os.Getenv("ARTEMIS_PASSWORD")))
}

func encodeGob(message any) ([]byte, error) {
	gob.Register(message)
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&message) // Pass pointer to interface so Encode sees a value of interface type.
	if err != nil {
		return nil, fmt.Errorf("could not encode as gob: %v", err)
	}
	return buff.Bytes(), nil
}

func decodeGob[T any](message []byte) (*T, error) {
	gob.Register(*new(T))
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var msg any
	err := dec.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("could not decode gob: %v", err)
	}
	m := msg.(T)
	return &m, nil
}
