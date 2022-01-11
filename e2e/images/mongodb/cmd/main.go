package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	connectStr string
	dataBase   string
	collection string
)

func main() {
	flag.StringVar(&connectStr, "connectStr", "", "")
	flag.StringVar(&dataBase, "dataBase", "", "")
	flag.StringVar(&collection, "collection", "", "")
	flag.Parse()

	ctx := context.Background()
	options := options.Client().ApplyURI(connectStr)
	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		fmt.Println("Failed connect with mongoDB:", err)
		log.Fatal(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("Failed disConnect with mongoDB: %v", err)
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			fmt.Printf("Failed disConnect with mongoDB: %v", err)
			log.Fatal(err)
		}
	}()

	err = client.Database(dataBase).Collection(collection).FindOneAndUpdate(ctx, bson.D{{"state", "running"}}, bson.M{"$set": bson.M{"state": "finished"}}).Err()
	if err != nil {
		fmt.Println(err)
	}
	log.Print("Work done,state changed.")
}
