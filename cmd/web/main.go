package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/asstronom/foodadvisor/pkg/app"
	"github.com/asstronom/foodadvisor/pkg/config"
	"github.com/asstronom/foodadvisor/pkg/db"
)

var (
	isDev bool
)

func parseFlags() {
	flag.BoolVar(&isDev, "dev", false, "set development mode")
	flag.Parse()
}

func main() {
	parseFlags()

	dbconf, err := config.ParseDBConnConfig(`configs\dbconf.yaml`)
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
	}

	if isDev {
		dbconf, err = config.ParseDBConnConfig(`configs\dbconf_dev.yaml`)
		if err != nil {
			log.Fatalf("error parsing dev conf: %s\n", err)
		}
	}
	fmt.Println(dbconf)
	log.Fatal()
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
