package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/Azure/azure-event-hubs-go/v3/eph"
	"github.com/Azure/azure-event-hubs-go/v3/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"
)

func main() {
	ehubNamespaceConnectionString := os.Getenv("EVENTHUB_CONNECTION_STRING")
	blobStorageConnectionString := os.Getenv("STORAGE_CONNECTION_STRING")
	blobContainerName := os.Getenv("CHECKPOINT_CONTAINER")

	// Azure Event Hub connection string
	parsed, err := conn.ParsedConnectionFromStr(ehubNamespaceConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create a new Azure Storage Leaser / Checkpointer
	storageAccountName, accountKey, err := parseAzureStorageConnectionString(blobStorageConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	cred, err := azblob.NewSharedKeyCredential(storageAccountName, accountKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	leaserCheckpointer, err := storage.NewStorageLeaserCheckpointer(cred, storageAccountName, blobContainerName, azure.PublicCloud)
	if err != nil {
		fmt.Println(err)
		return
	}

	// SAS token provider for Azure Event Hubs
	provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(parsed.KeyName, parsed.Key))
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// create a new EPH processor
	processor, err := eph.New(ctx, parsed.Namespace, parsed.HubName, provider, leaserCheckpointer, leaserCheckpointer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// register a message handler -- many can be registered
	handlerID, err := processor.RegisterHandler(ctx,
		func(c context.Context, e *eventhub.Event) error {
			fmt.Println(string(e.Data))
			return nil
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("handler id: %q is running\n", handlerID)

	// unregister a handler to stop that handler from receiving events
	// processor.UnregisterHandler(ctx, handleID)

	// start handling messages from all of the partitions balancing across multiple consumers
	err = processor.StartNonBlocking(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	err = processor.Close(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func parseAzureStorageConnectionString(connectionString string) (string, string, error) {
	parts := strings.Split(connectionString, ";")

	getValue := func(pair string) string {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			return parts[1]
		}
		return ""
	}

	var name, key string
	for _, v := range parts {
		switch {
		case strings.HasPrefix(v, "AccountName"):
			name = getValue(v)
		case strings.HasPrefix(v, "AccountKey"):
			key = getValue(v)

		}
	}

	if name == "" || key == "" {
		return "", "", errors.New("can't parse storage connection string. Missing key or name")
	}

	return name, key, nil
}
