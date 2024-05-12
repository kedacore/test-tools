package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

type msgMode string

const (
	PUT msgMode = "PUT"
	GET         = "GET"
)

type connectionConfig struct {
	host         string
	port         string
	queueManager string
	channel      string
	queueName    string

	appUsername string
	appPassword string
}

func setConnectionConfigFromEnv(connCfg *connectionConfig) error {
	lookupRequiredEnvKey := func(envKey string) (string, error) {
		val, ok := os.LookupEnv(envKey)
		if !ok {
			return "", fmt.Errorf("missing value for the required environment variable '%s'", envKey)
		}
		return val, nil
	}

	host, err := lookupRequiredEnvKey("HOST")
	if err != nil {
		return err
	}
	connCfg.host = host

	port, err := lookupRequiredEnvKey("PORT")
	if err != nil {
		return err
	}
	connCfg.port = port

	queueManager, err := lookupRequiredEnvKey("QUEUE_MANAGER")
	if err != nil {
		return err
	}
	connCfg.queueManager = queueManager

	channel, err := lookupRequiredEnvKey("CHANNEL")
	if err != nil {
		return err
	}
	connCfg.channel = channel

	queueName, err := lookupRequiredEnvKey("QUEUE_NAME")
	if err != nil {
		return err
	}
	connCfg.queueName = queueName

	connCfg.appUsername = os.Getenv("APP_USERNAME")
	connCfg.appPassword = os.Getenv("APP_PASSWORD")

	return nil
}

func createConnection(connCfg connectionConfig) (ibmmq.MQQueueManager, error) {
	log.Println("setting up connection to MQ")
	// Allocate the default MQ Connection Options (MQCNO) structure needed for the CONNX call.
	conn := ibmmq.NewMQCNO()

	if connCfg.appUsername != "" {
		log.Printf("username '%s' has been specified\n", connCfg.appUsername)
		// The MQ Security Parameters (MQCSP)
		csp := ibmmq.NewMQCSP()
		csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
		csp.UserId = connCfg.appUsername
		csp.Password = connCfg.appPassword

		conn.SecurityParms = csp
	}

	// Allocate the default MQ Channel Definition (MQCD) with the required fields
	cd := ibmmq.NewMQCD()
	cd.ChannelName = connCfg.channel
	cd.ConnectionName = connCfg.host + "(" + connCfg.port + ")"
	log.Printf("connecting to %s\n", cd.ConnectionName)

	conn.ClientConn = cd
	conn.Options = ibmmq.MQCNO_CLIENT_BINDING
	log.Printf("attempting connection to queue manager(QMGR): %s\n", connCfg.queueManager)

	return ibmmq.Connx(connCfg.queueManager, conn)
}

func openQueue(qMgrObject ibmmq.MQQueueManager, queueName string, mode msgMode) (ibmmq.MQObject, error) {
	var (
		qObject ibmmq.MQObject
		err     error
	)

	mqod := ibmmq.NewMQOD()
	openOptions := ibmmq.MQOO_OUTPUT

	switch mode {
	case PUT:
		mqod.ObjectType = ibmmq.MQOT_Q
		mqod.ObjectName = queueName

	case GET:
		openOptions = ibmmq.MQOO_INPUT_SHARED
		mqod.ObjectType = ibmmq.MQOT_Q
		mqod.ObjectName = queueName

	default:
		return qObject, fmt.Errorf("invalid msgMode, it must be either 'PUT' or 'GET'")
	}

	qObject, err = qMgrObject.Open(mqod, openOptions)

	return qObject, err
}
