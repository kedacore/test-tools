package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func main() {
	var connCfg connectionConfig
	if err := setConnectionConfigFromEnv(&connCfg); err != nil {
		log.Fatal(err)
	}

	qMgr, err := createConnection(connCfg)
	if err != nil {
		log.Fatal("unable to establish connection to the queue manager(QMGR): ", err)
	}
	defer func() {
		if err := qMgr.Disc(); err != nil {
			log.Println("error while disconnecting from the queue manager(QMGR): ", err)
		}
	}()

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	switch mode {
	case "producer":
		if err := produceMessages(qMgr, connCfg.queueName); err != nil {
			log.Fatal(err)
		}
	case "consumer":
		if err := consumeMessages(qMgr, connCfg.queueName); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid mode '%s', it must be either 'producer' or 'consumer'\n", mode)
	}
}

const (
	defaultProducerConsumerSleepTime = 1
	defaultProducedNumMessages       = 100
)

type message struct {
	Content string `json:"content"`
	Number  int    `json:"number"`
}

func consumeMessages(qMgr ibmmq.MQQueueManager, queueName string) error {
	var err error

	qObject, err := openQueue(qMgr, queueName, GET)
	if err != nil {
		return fmt.Errorf("failed to open the queue %s: %w", queueName, err)
	}
	defer func() {
		if cerr := qObject.Close(0); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Sleep to demonstrate scaling
	sleepTimeString := os.Getenv("CONSUMER_SLEEP_TIME")
	sleepTime, err := strconv.Atoi(sleepTimeString)
	if err != nil {
		sleepTime = defaultProducerConsumerSleepTime
	}

	for {
		// The Get requires control structures, the Message Descriptor (MQMD)
		// and Get Options (MQGMO). We create those with default values.
		getmqmd := ibmmq.NewMQMD()
		gmo := ibmmq.NewMQGMO()
		// The default options are OK, but it's always a good idea to be explicit
		// about transactional boundaries as not all platforms behave the same way.
		gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT
		// Set options to wait for a maximum of 3 seconds for any new message to arrive
		gmo.Options |= ibmmq.MQGMO_WAIT
		// Unlimited get
		gmo.WaitInterval = ibmmq.MQWI_UNLIMITED

		buf := make([]byte, 1024)
		datalen, err := qObject.Get(getmqmd, gmo, buf)
		if err != nil {
			mqret := err.(*ibmmq.MQReturn)
			if mqret.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
				log.Println("no messages are currently available to consume, this will stop the consumer loop")
				// No message is available is an expected situation, and so we don't handle it as a real error.
				err = nil
			}
			break
		}

		var msg message
		err = json.Unmarshal(buf[:datalen], &msg)
		if err != nil {
			return fmt.Errorf("failed to parse JSON message: %w", err)
		}
		log.Printf("received message from %s with content: %s, number: %d\n", queueName, msg.Content, msg.Number)

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	return err
}

func produceMessages(qMgr ibmmq.MQQueueManager, queueName string) error {
	var err error

	qObject, err := openQueue(qMgr, queueName, PUT)
	if err != nil {
		return fmt.Errorf("failed to open the queue %s: %w", queueName, err)
	}
	defer func() {
		if cerr := qObject.Close(0); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Sleep to demonstrate scaling
	sleepTimeString := os.Getenv("PRODUCER_SLEEP_TIME")
	sleepTime, err := strconv.Atoi(sleepTimeString)
	if err != nil {
		sleepTime = defaultProducerConsumerSleepTime
	}

	numMessages, err := strconv.Atoi(os.Getenv("NUM_MESSAGES"))
	if err != nil {
		numMessages = defaultProducedNumMessages
	}

	for i := 0; i < numMessages; i++ {
		// The PUT requires control structures, the Message Descriptor (MQMD)
		// and Put Options (MQPMO). Create those with default values.
		putmqmd := ibmmq.NewMQMD()
		pmo := ibmmq.NewMQPMO()
		// The default options are OK, but it's always
		// a good idea to be explicit about transactional boundaries as
		// not all platforms behave the same way.
		pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
		// Tell MQ what the message body format is. In this case, a text string
		putmqmd.Format = ibmmq.MQFMT_STRING

		msg := &message{
			Content: "Msg created at " + time.Now().Format(time.RFC3339),
			Number:  i + 1,
		}
		buf, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshalling the message data to send: %w", err)
		}

		err = qObject.Put(putmqmd, pmo, buf)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		log.Printf("message: %s sent to %s\n", buf, strings.TrimSpace(qObject.Name))

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	return err
}
