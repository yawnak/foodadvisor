package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/asstronom/foodadvisor/cmd/web/config"
	"github.com/asstronom/foodadvisor/pkg/db"
	migrations "github.com/asstronom/foodadvisor/pkg/db/migrations"
	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/sethvargo/go-envconfig"
)

var (
	conf = map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"USER":     "user",
		"PASSWORD": "mypassword",
		"NAME":     "food",
	}
)

func OpenDB() *db.FoodDB {
	ml := envconfig.MapLookuper(conf)
	dbconf := config.DBConnConfig{}
	err := envconfig.ProcessWith(context.Background(), &dbconf, ml)
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
	return foodRepo
}

func TestUserCRUD(t *testing.T) {
	foodRepo := OpenDB()
	user := domain.User{
		Username:       "pudgbooster",
		Password:       "hater",
		ExpirationDays: 5,
	}
	id, err := foodRepo.CreateUser(context.Background(), &user)
	if err != nil {
		t.Fatalf("error creating user: %s\n", err)
	}
	t.Logf("created userid: %d\n", id)

	foundUser, err := foodRepo.GetUserById(context.Background(), id)

	if err != nil {
		t.Fatalf("error getting user: %s\n", err)
	}

	t.Logf("got user: %v", foundUser)

	foundUser.Username = "lucky"

	err = foodRepo.UpdateUser(context.Background(), foundUser)
	if err != nil {
		t.Fatalf("error updating user: %s\n", err)
	}

	foundUser, err = foodRepo.GetUserById(context.Background(), id)

	if err != nil {
		t.Fatalf("error getting updated user: %s\n", err)
	}

	t.Logf("got user: %v", foundUser)

	err = foodRepo.DeleteUser(context.Background(), id)
	if err != nil {
		t.Fatalf("error deleting user: %s\n", err)
	}
}
