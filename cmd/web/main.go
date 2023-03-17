package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asstronom/foodadvisor/cmd/web/config"
	"github.com/asstronom/foodadvisor/pkg/app"
	"github.com/asstronom/foodadvisor/pkg/db"
)

func main() {
	dbconf, err := config.ParseDBConnConfig(`C:\Users\danya\Documents\foodadvisor\configs\dbconf.yaml`)
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
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

	advisor, _ := app.NewFoodAdvisor(foodRepo)
	fmt.Println(advisor)
}
