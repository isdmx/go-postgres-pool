package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/isdmx/go-postgres-pool/config"
	"github.com/isdmx/go-postgres-pool/db"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	ctx := context.Background()
	dbPool, err := db.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Could not create database pool: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	row, _ := dbPool.QueryContext(ctx, "SELECT * FROM generate_series(2,50);")

	for row.Next() {
		var value int
		if err := row.Scan(&value); err != nil {
			log.Fatalf("Could not scan row: %v", err)
		}
		log.Printf("Value: %d", value)
		time.Sleep(5 * time.Second)
	}

	fmt.Println("Successfully connected to the database!")
}
