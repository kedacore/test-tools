package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	url := os.Args[1]
	queueName := "hello"
	interval := 0

	messageCount, err := strconv.Atoi(os.Args[2])
	failOnError(err, "Failed to parse second arg as messageCount : int")

	if len(os.Args) >= 4 {
		queueName = os.Args[3]
	}

	if len(os.Args) == 5 {
		interval, err = strconv.Atoi(os.Args[4])
	}
	failOnError(err, "Failed to parse fourth arg as messages publishing interval (in seconds) : int")

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if interval > 0 {
		log.Printf(" [*] Publishing %d message(s) per %d second(s). To exit press CTRL+C", messageCount, interval)
	}

	for {
		for i := 0; i < messageCount; i++ {
			body := fmt.Sprintf("Hello World: %d", i)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			log.Printf(" [x] Sent %s", body)
			failOnError(err, "Failed to publish a message")
		}

		if interval == 0 {
			break
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
