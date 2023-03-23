package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/yawnak/foodadvisor/pkg/domain"
)

const (
	host     string = "localhost"
	username string = "user"
	password string = "mypassword"
	port     string = "5432"
	dbname   string = "food"
)

func TestCreateRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()
	err = foodRepo.CreateRole(context.Background(), domain.Role{
		Name:        "admin",
		Permissions: map[domain.Permission]struct{}{domain.PermEditUserRole: {}}})
	if err != nil {
		t.Errorf("error creating role: %s", err)
	}
}

func TestGetRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()
	role, err := foodRepo.GetRole(context.Background(), "admin")
	if err != nil {
		t.Errorf("error getting role: %s", err)
	}
	fmt.Println(role)
}
