package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/yawnak/foodadvisor/internal/domain"
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
	err = foodRepo.CreateRole(context.Background(),
		domain.NewRole("owner", domain.PermEditRoles, domain.PermEditUserRole))
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

func TestUpdateRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()
	err = foodRepo.UpdateRole(context.Background(),
		&domain.Role{
			Name: "pudge",
			Permissions: map[domain.Permission]struct{}{
				domain.PermEditRoles:    {},
				domain.PermEditUserRole: {},
			},
		})
	if err != nil {
		t.Errorf("error getting role: %s", err)
	}
}

func TestDeleteRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()

	err = foodRepo.DeleteRole(context.Background(), "owner")
	if err != nil {
		t.Errorf("error deleting role: %s", err)
	}
}

func TestGetUserRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()

	role, err := foodRepo.GetUserRole(context.Background(), 1)
	if err != nil {
		t.Errorf("error getting user role: %s", err)
	}
	fmt.Println(role)
}

func TestUpdateUserRole(t *testing.T) {
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	foodRepo, err := Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer foodRepo.Close()

	err = foodRepo.UpdateUserRole(context.Background(), 2, `user`)
	if err != nil {
		t.Errorf("error updating user role: %s", err)
	}
}
