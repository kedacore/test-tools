package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	mode := flag.String("mode", "", "mssql tool mode. Required. Allowed values = [producer, consumer]")
	flag.Parse()
	switch *mode {
	case "producer":
		runProducer()
	case "consumer":
		runConsumer()
	default:
		log.Fatalf("Unsupported mode. %s", *mode)
	}
}

func runProducer() {
	log.Println("mode = Producer")
	itemCount := 10
	connString := getConnString()
	log.Printf("Inserting %d items into the 'tasks' table...", itemCount)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed: ", err.Error())
	}
	defer conn.Close()
	for range itemCount {
		_, err := conn.Exec("INSERT INTO tasks ([status]) VALUES ('queued')")
		if err != nil {
			log.Fatal("Insert task failed: ", err.Error())
		}
	}

	log.Printf("Inserting %d records succesfully!", itemCount)
}

func runConsumer() {
	log.Println("mode = Consumer")
	connString := getConnString()

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed: ", err.Error())
	}
	defer conn.Close()
	taskCmd := "UPDATE TOP (1) tasks SET [status] = 'running' OUTPUT inserted.[id] WHERE [status] = 'queued'"
	var taskId int
	for {
		err := conn.QueryRow(taskCmd).Scan(&taskId)
		if err == sql.ErrNoRows {
			log.Printf("No queue task at the moment...")
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			log.Fatal("Query queue task failed: ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		// Simulate work
		log.Printf("Simulate work for taskId: %d ....", taskId)
		time.Sleep(30 * time.Second)
		_, err = conn.Exec("UPDATE tasks SET [status] = 'complete' WHERE [id] = ?", taskId)
		if err != nil {
			log.Fatal("Update completed task failed: ", err.Error())
		}
	}
}

func getConnString() string {
	return os.Getenv("SQL_CONNECTION_STRING")
}
