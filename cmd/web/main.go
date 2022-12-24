package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asstronom/foodadvisor/cmd/web/config"
	"github.com/asstronom/foodadvisor/pkg/db"
	migrations "github.com/asstronom/foodadvisor/pkg/db/migrations"
	"github.com/asstronom/foodadvisor/pkg/ui"
)

func main() {
	// dbconf, err := config.ParseDBConnConfigEnv(context.Background(), "DB_")
	// if err != nil {
	// 	log.Fatalf("error parsing db conf: %s\n", err)
	// }

	dbconf, err := config.ParseDBConnConfig(`C:\Users\danya\Documents\foodadvisor\configs\dbconf.yaml`)
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
	}

	migrateurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	pathToMigrations := "../pkg/db/migrations/sql" //"../../pkg/db/migrations/sql"
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
	defer foodRepo.Close()
	fmt.Println("Hello world!")

	cli := ui.UICli{}
	err = cli.Run()
	if err != nil {
		log.Fatalf("error in ui: %s", err)
	}
}
