package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	taskInstanceCount, err := strconv.Atoi(os.Getenv("TASK_INSTANCES_COUNT"))
	if err != nil {
		panic("Error parsing TASK_INSTANCES_COUNT")
	}

	actionType := os.Args[1]
	if actionType == "insert" {
		addTaskInstances(taskInstanceCount)
	}
	if actionType == "update" {
		updateTaskInstances(taskInstanceCount)
		// keep the worker running indefinitely
		exit := make(chan string)
		for {
			select {
			case <-exit:
				os.Exit(0)
			}
		}
	}

}

type taskInstanceState string

const (
	queued   taskInstanceState = "queued"
	procesed taskInstanceState = "processed"
)

func addTaskInstances(numberOfTaskInstancesToAdd int) {
	log.Printf("Adding %v task instances", numberOfTaskInstancesToAdd)
	db := getDB()
	defer db.Close()
	for i := 0; i < numberOfTaskInstancesToAdd; i++ {
		log.Printf("Inserting %v of %v", (i + 1), numberOfTaskInstancesToAdd)
		insert, err := db.Query("INSERT INTO task_instance (state) VALUES ($1)", queued)
		insert.Close()
		if err != nil {
			panic(err.Error())
		}
	}
}

type taskInstance struct {
	id    int
	state string
}

func updateTaskInstances(numberOfTaskInstancesToUpdate int) {
	log.Printf("Updating %v task instances", numberOfTaskInstancesToUpdate)
	db := getDB()
	defer db.Close()
	for i := 0; i < numberOfTaskInstancesToUpdate; i++ {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
		if err != nil {
			log.Fatalf("Error creating transaction %v", err)
		}

		var taskInstance taskInstance
		err = tx.QueryRowContext(ctx, "SELECT id, state FROM task_instance WHERE state = $1 LIMIT 1", queued).Scan(&taskInstance.id, &taskInstance.state)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
			time.Sleep(2 * time.Second)
		}
		log.Printf("Updating task instance id %v", taskInstance.id)
		updateStmt, err := tx.PrepareContext(ctx, "UPDATE task_instance SET state = $1 WHERE id = $2")
		if err != nil {
			tx.Rollback()
			log.Fatalf("Rolling back transaction %v", err)
			panic(err.Error())
		}
		_, err = updateStmt.ExecContext(ctx, procesed, taskInstance.id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Rolling back transaction %v", err)
		}
		updateStmt.Close()
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Error committing transaction! Error is %v", err)
		} else {
			log.Printf("Transaction committed for id %v", taskInstance.id)
		}
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	db.SetMaxOpenConns(5)
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	return db
}
