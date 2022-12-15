package main

import (
	"fmt"
	"log"

	"github.com/asstronom/foodadvisor/cmd/web/config"
	migrations "github.com/asstronom/foodadvisor/pkg/db/migrations"
)

func main() {
	dbconf, err := config.ParseDBConnConfigEnv()
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
	}

	migrateurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbconf.DB_USER, dbconf.DB_PASSWORD, dbconf.DB_HOST, dbconf.DB_PORT, dbconf.DB_NAME)
	pathToMigrations := "../../pkg/db/migrations/sql"
	err = migrations.MigrateUp(pathToMigrations, migrateurl)
	if err != nil {
		log.Fatalf("error migrating db: %s\n", err)
	}

	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbconf.DB_USER, dbconf.DB_PASSWORD, dbconf.DB_HOST, dbconf.DB_PORT, dbconf.DB_NAME)

	fmt.Println(dburl)
	fmt.Println("Hello world!")
}
