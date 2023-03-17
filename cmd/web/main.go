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
	//flags
	isDev bool
)

// parse all flags
func parseFlags() {
	flag.BoolVar(&isDev, "dev", false, "set development mode")
	flag.Parse()
}

func main() {
	parseFlags()

	//parsing configs
	var pathToDBConf string
	if !isDev {
		pathToDBConf = `configs\dbconf.yaml`
	} else {
		pathToDBConf = `configs\dbconf_dev.yaml`
	}

	dbconf, err := config.ParseDBConnConfig(pathToDBConf)
	if err != nil {
		log.Fatalf("error parsing db conf: %s\n", err)
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
