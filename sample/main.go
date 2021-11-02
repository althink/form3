package main

import (
	"context"
	"log"

	"github.com/althink/form3"
	"github.com/althink/form3/accounts"
)

func main() {
	ctx := context.Background()
	orgID := "0de1f73f-8af2-4316-86f9-325ce9755cb6"

	f3, err := form3.NewClient()
	if err != nil {
		log.Fatal("Failed to create client", err)
	}

	acc, err := f3.Accounts.Create(ctx, accounts.NewWithGenID(orgID, &accounts.Attributes{
		Country: "PL",
		Name:    []string{"John Smith"},
	}))
	if err != nil {
		log.Fatal("Failed to create account", err)
	}
	log.Printf("Account %s created successfully\n", acc.Data.ID)

	_, err = f3.Accounts.Fetch(ctx, acc.Data.ID)
	if err != nil {
		log.Fatal("Failed to fetch account", err)
	}
	log.Printf("Account %s fetched successfully\n", acc.Data.ID)

	err = f3.Accounts.Delete(ctx, acc.Data.ID, *acc.Data.Version)
	if err != nil {
		log.Fatal("Failed to delete account", err)
	}
	log.Printf("Account %s deleted successfully\n", acc.Data.ID)
}
