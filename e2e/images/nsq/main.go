package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

type Handler struct{}

func (h *Handler) HandleMessage(m *nsq.Message) error {
	log.Printf("Received message: %s", m.Body)
	return nil
}

func nsqConsumer(config *nsq.Config, nsqlookupdHTTPAddress, topic, channel string) error {
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return err
	}

	consumer.AddHandler(&Handler{})

	err = consumer.ConnectToNSQLookupd(nsqlookupdHTTPAddress)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()

	return nil
}

func nsqProducer(config *nsq.Config, nsqdTCPAddress string, topic string, messageCount int) error {
	producer, err := nsq.NewProducer(nsqdTCPAddress, config)
	if err != nil {
		return err
	}

	responseChan := make(chan *nsq.ProducerTransaction, messageCount)
	for i := 0; i < messageCount; i++ {
		err := producer.PublishAsync(topic, []byte(fmt.Sprintf("%d", i)), responseChan)
		if err != nil {
			return err
		}
	}

	for i := 0; i < messageCount; i++ {
		trans := <-responseChan
		if trans.Error != nil {
			return trans.Error
		}
	}

	producer.Stop()

	return nil
}

func main() {
	mode := flag.String("mode", "", "consumer or producer")
	topic := flag.String("topic", "", "topic name")
	channel := flag.String("channel", "", "channel name")
	nsqlookupdHTTPAddress := flag.String("nsqlookupd-http-address", "", "nsqlookupd HTTP address")
	messageCount := flag.Int("message-count", 1, "number of messages to send")
	nsqdTCPAddress := flag.String("nsqd-tcp-address", "", "nsqd TCP address")
	flag.Parse()

	config := nsq.NewConfig()

	switch *mode {
	case "consumer":
		log.Println("Consumer mode")
		if *topic == "" || *channel == "" || *nsqlookupdHTTPAddress == "" {
			log.Fatalf("topic, channel, and nsqlookupd-http-address are required\n")
		}
		if err := nsqConsumer(config, *nsqlookupdHTTPAddress, *topic, *channel); err != nil {
			log.Fatalf("read from nsq failed: %w\n", err)
		}
	case "producer":
		log.Println("Producer mode")
		if *topic == "" || *nsqdTCPAddress == "" {
			log.Fatalf("topic and nsqd-tcp-address are required\n")
		}
		if err := nsqProducer(config, *nsqdTCPAddress, *topic, *messageCount); err != nil {
			log.Fatalf("write to nsq failed: %w\n", err)
		}
	default:
		log.Fatalf("unknown mode: %s\n", *mode)
	}
}
