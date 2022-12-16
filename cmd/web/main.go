package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asstronom/foodadvisor/cmd/web/config"
	"github.com/asstronom/foodadvisor/pkg/db"
	migrations "github.com/asstronom/foodadvisor/pkg/db/migrations"
)

func main() {
	dbconf, err := config.ParseDBConnConfigEnv(context.Background(), "DB_")
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
	}

	migrateurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	pathToMigrations := "../../pkg/db/migrations/sql"
	err = migrations.MigrateUp(pathToMigrations, migrateurl)
	if err != nil {
		log.Fatalf("error migrating db: %s\n", err)
	}

	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	fmt.Println(dburl)
	foodRepo, err := db.Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	fmt.Println("Hello world!")
	foodRepo.Close()
}
