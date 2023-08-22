package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yawnak/foodadvisor/pkg/app"
	"github.com/yawnak/foodadvisor/pkg/config"
	"github.com/yawnak/foodadvisor/pkg/db"
	"github.com/yawnak/foodadvisor/pkg/server"
)

var (
	//flags
	wenv string
)

// parse all flags
func parseFlags() {
	flag.StringVar(&wenv, "wenv", "prod", "set environment mode")
	flag.Parse()
}

func main() {
	var err error
	parseFlags()

	var dbconf config.DBConnConfig
	log.Println("environment:", wenv)
	err = config.BindConfig(&dbconf, "configs/conf.yaml", fmt.Sprintf("configs/conf.%s.yaml", wenv))
	if err != nil {
		log.Fatalln(err)
	}

	dbconf.Password = os.Getenv("DB_PASSWORD")

	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	foodRepo, err := db.Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()
	log.Println("successfully connected to food database")

	advisor, _ := app.NewFoodAdvisor(foodRepo)
	if err != nil {
		log.Fatalf("error creating advisor service: %s", err)
	}

	srv, err := server.NewServer(advisor)
	if err != nil {
		log.Fatalf("error creating server: %s", err)
	}
	log.Fatalln("error listening and serving: ", srv.ListenAndServe("8080"))
}
