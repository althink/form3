//go:build integration

package form3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/althink/form3/accounts"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.M) {
	err := waitUntilHealthy(60 * time.Second)
	if err != nil {
		log.Fatal("Skipping tests: ", err)
	}

	t.Run()
}

func TestIntegration_Accounts_CreateFetchAndDeleteSuccess(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	orgID := uuid.New().String()
	acc := accounts.NewWithGenID(orgID,
		&accounts.Attributes{
			Country:                 "GB",
			BaseCurrency:            "GBP",
			AccountNumber:           "41426819",
			BankID:                  "400300",
			BankIDCode:              "GBDSC",
			Bic:                     "NWBKGB22",
			Iban:                    "GB11NWBK40030041426819",
			Name:                    []string{"Samantha Holder"},
			AlternativeNames:        []string{"Sam Holder"},
			AccountClassification:   strPtr("Personal"),
			JointAccount:            boolPtr(false),
			AccountMatchingOptOut:   boolPtr(false),
			SecondaryIdentification: strPtr("A1B2C3D4"),
			Switched:                boolPtr(false),
			Status:                  strPtr("confirmed"),
		})

	// when
	created, err := c.Accounts.Create(context.Background(), acc)

	// then
	require.Empty(t, err)
	require.Equal(t, acc.ID, created.Data.ID, "Invalid ID")
	require.Equal(t, acc.OrganisationID, created.Data.OrganisationID, "Invalid OrganisationID")
	require.Equal(t, acc.Type, created.Data.Type, "Invalid Type")
	require.Equal(t, int64(0), *created.Data.Version, "Invalid Version")
	require.NotEmpty(t, created.Links.Self)
	if diff := cmp.Diff(acc.Attributes, created.Data.Attributes); diff != "" {
		t.Errorf("Created account mismatch (-want +got):\n%s", diff)
	}

	// when
	fetched, err := c.Accounts.Fetch(ctx, created.Data.ID)

	// then
	require.Empty(t, err)
	require.Equal(t, acc.ID, fetched.Data.ID, "Invalid ID")
	require.Equal(t, acc.OrganisationID, fetched.Data.OrganisationID, "Invalid OrganisationID")
	require.Equal(t, acc.Type, fetched.Data.Type, "Invalid Type")
	require.Equal(t, int64(0), *fetched.Data.Version, "Invalid Version")
	require.NotEmpty(t, created.Links.Self)
	if diff := cmp.Diff(acc.Attributes, fetched.Data.Attributes); diff != "" {
		t.Errorf("Fetched account mismatch (-want +got):\n%s", diff)
	}

	// when
	err = c.Accounts.Delete(ctx, fetched.Data.ID, *fetched.Data.Version)

	// then
	require.Empty(t, err)

	// when
	_, err = c.Accounts.Fetch(ctx, created.Data.ID)

	// then
	require.IsType(t, &accounts.AccountNotFoundError{}, err, "Invalid error type")
}

func TestIntegration_Accounts_CreateFailed_InvalidData(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()

	// when
	_, err := c.Accounts.Create(ctx, &accounts.Data{Attributes: &accounts.Attributes{}})

	// then
	require.IsType(t, &accounts.InvalidDataError{}, err, "Invalid error type")
	require.NotEmpty(t, err.(*accounts.InvalidDataError).Msg)
}

func TestIntegration_Accounts_CreateFailed_AlreadyExists(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := uuid.New().String()
	orgID := uuid.New().String()
	accountExists(t, c, accID)

	// when
	_, err := c.Accounts.Create(ctx, accounts.New(accID, orgID, &accounts.Attributes{
		Country: "GB",
		Name:    []string{"Samantha Holder"},
	}))

	// then
	require.IsType(t, &accounts.AccountAlreadyExistsError{}, err, "Invalid error type")
	require.Equal(t, accID, err.(*accounts.AccountAlreadyExistsError).ID)
}

func TestIntegration_Accounts_FetchFailed_InvalidData(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := "invalid"

	// when
	_, err := c.Accounts.Fetch(ctx, accID)

	// then
	require.IsType(t, &accounts.InvalidDataError{}, err, "Invalid error type")
	require.NotEmpty(t, err.(*accounts.InvalidDataError).Msg)
}

func TestIntegration_Accounts_FetchFailed_UnknownAccount(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := "6acb52e8-375b-453e-a3a9-9110c9aca283"

	// when
	_, err := c.Accounts.Fetch(ctx, accID)

	// then
	require.IsType(t, &accounts.AccountNotFoundError{}, err, "Invalid error type")
	require.Equal(t, accID, err.(*accounts.AccountNotFoundError).ID)
}

func TestIntegration_Accounts_DeleteFailed_UnknownAccount(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := "6acb52e8-375b-453e-a3a9-9110c9aca283"

	// when
	err := c.Accounts.Delete(ctx, accID, 0)

	// then
	require.IsType(t, &accounts.AccountNotFoundError{}, err, "Invalid error type")
	require.Equal(t, accID, err.(*accounts.AccountNotFoundError).ID)
}

func TestIntegration_Accounts_DeleteFailed_InvalidVersion(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := uuid.New().String()
	invalidVer := int64(7)
	accountExists(t, c, accID)

	// when
	err := c.Accounts.Delete(ctx, accID, invalidVer)

	// then
	require.IsType(t, &accounts.InvalidVersionError{}, err, "Invalid error type")
	require.Equal(t, invalidVer, err.(*accounts.InvalidVersionError).Ver)
}

func TestIntegration_Accounts_DeleteFailed_InvalidData(t *testing.T) {
	// given
	c := setUpClient()
	ctx := context.Background()
	accID := "invalid"

	// when
	err := c.Accounts.Delete(ctx, accID, 0)

	// then
	require.IsType(t, &accounts.InvalidDataError{}, err, "Invalid error type")
	require.NotEmpty(t, err.(*accounts.InvalidDataError).Msg)
}

func setUpClient() *Form3 {
	u, err := url.Parse(resolveApiUrl())
	if err != nil {
		log.Fatal(err)
	}
	f3, err := NewClient(WithBaseURL(*u))
	if err != nil {
		log.Fatal(err)
	}
	return f3
}

func resolveApiUrl() string {
	host := os.Getenv("FORM3_HOST")
	if host == "" {
		host = "http://localhost:8080/v1/"
	}
	return host
}

func accountExists(t *testing.T, c *Form3, id string) {
	_, err := c.Accounts.Create(context.Background(), accounts.New(id, uuid.New().String(),
		&accounts.Attributes{
			Country: "PL",
			Name:    []string{"John Smith"},
		}))
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}
}

func waitUntilHealthy(timeout time.Duration) error {
	baseUrl := resolveApiUrl()
	c := make(chan string, 1)
	go func() {
		for {
			resp, err := http.Get(baseUrl + "health")
			if err == nil && resp.StatusCode == 200 {
				c <- "success"
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	select {
	case <-c:
		// health endpoint returned 200
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Service not ready in %s", timeout)
	}
}

func strPtr(val string) *string {
	return &val
}

func boolPtr(val bool) *bool {
	return &val
}
