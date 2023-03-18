package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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
	var err error
	parseFlags()

	var dbconf config.DBConnConfig
	var env string
	if isDev {
		env = "dev"
	} else {
		env = "prod"
	}
	err = config.BindConfig(&dbconf, "configs/conf.yaml", fmt.Sprintf("configs/conf.%s.yaml", env))
	if err != nil {
		log.Fatalln(err)
	}

	dbconf.Password = os.Getenv("DB_PASSWORD")

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
	log.Println("successfully connected to food database")

	advisor, _ := app.NewFoodAdvisor(foodRepo)
	fmt.Println(advisor)
}
