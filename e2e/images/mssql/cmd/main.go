package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

var sqlIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func main() {
	mode := flag.String("mode", "", "mssql tool mode. Required. Allowed values = [producer, consumer]")
	table := flag.String("table", "tasks", "table used to produce or consume tasks")
	flag.Parse()

	tableName, err := quoteIdentifier(*table)
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "producer":
		runProducer(tableName)
	case "consumer":
		runConsumer(tableName)
	default:
		log.Fatalf("Unsupported mode. %s", *mode)
	}
}

func runProducer(tableName string) {
	log.Println("mode = Producer")
	itemCount := 10
	connString := getConnString()
	log.Printf("Inserting %d items into the %s table...", itemCount, tableName)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed: ", err.Error())
	}
	defer conn.Close()
	insertCmd := fmt.Sprintf("INSERT INTO %s ([status]) VALUES ('queued')", tableName)
	for range itemCount {
		_, err := conn.Exec(insertCmd)
		if err != nil {
			log.Fatal("Insert task failed: ", err.Error())
		}
	}

	log.Printf("Inserting %d records succesfully!", itemCount)
}

func runConsumer(tableName string) {
	log.Println("mode = Consumer")
	connString := getConnString()

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed: ", err.Error())
	}
	defer conn.Close()
	taskCmd := fmt.Sprintf("UPDATE TOP (1) %s SET [status] = 'running' OUTPUT inserted.[id] WHERE [status] = 'queued'", tableName)
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
		_, err = conn.Exec(fmt.Sprintf("UPDATE %s SET [status] = 'complete' WHERE [id] = ?", tableName), taskId)
		if err != nil {
			log.Fatal("Update completed task failed: ", err.Error())
		}
	}
}

func getConnString() string {
	return os.Getenv("SQL_CONNECTION_STRING")
}

func quoteIdentifier(identifier string) (string, error) {
	if !sqlIdentifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid table name %q: expected letters, digits, or underscores, starting with a letter or underscore", identifier)
	}
	return "[" + identifier + "]", nil
}
