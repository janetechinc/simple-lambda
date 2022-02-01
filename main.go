package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

type Request struct {
	ID int `json:"id"`
}

func setup(ctx context.Context) error {
	var err error

	conn, err = pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	_, err = conn.Exec(
		ctx,
		"CREATE TABLE IF NOT EXISTS counters"+
			" (id integer PRIMARY KEY, value integer)",
	)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func Handle(ctx context.Context, request Request) (int, error) {
	var value int
	err := conn.QueryRow(
		ctx,
		"INSERT INTO counters (id, value) VALUES ($1, 1)"+
			" ON CONFLICT (id) DO UPDATE SET value = counters.value + 1"+
			" RETURNING value",
		request.ID,
	).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("failed to query a counter: %w", err)
	}

	log.Printf("Counter: %d", value)
	return value, nil
}

func main() {
	err := setup(context.Background())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			log.Printf("Error closing database connection: %s", err)
		}
	}()

	isLambda := len(os.Getenv("_LAMBDA_SERVER_PORT")) > 0
	if isLambda {
		lambda.Start(Handle)
	} else {
		counter := 1
		flag.IntVar(&counter, "counter", 0, "Counter ID")
		flag.Parse()

		_, err := Handle(context.Background(), Request{ID: counter})
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
}
