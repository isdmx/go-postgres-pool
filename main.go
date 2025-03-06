package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/sourcegraph/conc/pool"

	"github.com/isdmx/go-postgres-pool/config"
	"github.com/isdmx/go-postgres-pool/db"
)

func main() {
	cfgW, err := config.LoadConfig("config.yaml", config.WithWrite())
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	ctx := context.Background()
	dbPoolW, err := db.New(ctx, cfgW)
	if err != nil {
		log.Fatalf("Could not create database pool: %v", err)
	}
	defer dbPoolW.Close()

	if err := dbPoolW.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	cfgR, err := config.LoadConfig("config.yaml", config.WithReadOnly())
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	dbPoolR, err := db.New(ctx, cfgR)
	if err != nil {
		log.Fatalf("Could not create database pool: %v", err)
	}
	defer dbPoolR.Close()

	if err := dbPoolW.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	cluster := []*sql.DB{dbPoolW, dbPoolR}

	p := pool.New().WithMaxGoroutines(200)

	for i := range 1000 {
		v := i
		p.Go(func() {
			db := cluster[v%len(cluster)]
			res, err := db.QueryContext(ctx, "select pg_sleep(4);")
			defer func() { _ = res.Close() }()

			if err != nil {
				log.Printf("Could not execute query: %v\n", err)
			}

			log.Printf("Executed query %d\n", v)
		})
	}

	fmt.Println("Successfully connected to the database!")
}
