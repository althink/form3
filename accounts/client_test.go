package accounts

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Accounts_CreateSuccess(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 201
	responseBody := `{
						"data": {
							"type": "accounts",
							"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
							"version": 0,
							"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
							"attributes": {
								"country": "GB",
								"base_currency": "GBP",
								"account_number": "41426819",
								"bank_id": "400300",
								"bank_id_code": "GBDSC",
								"bic": "NWBKGB22",
								"iban": "GB11NWBK40030041426819",
								"name": [
									"Samantha Holder"
								],
								"alternative_names": [
									"Sam Holder"
								],
								"status": "confirmed",
								"user_defined_data": [
								{
									"key": "Some account related key",
									"value": "Some account related value"
								}
								],
								"validation_type": "card",
								"reference_mask": "############",
								"acceptance_qualifier": "same_day",
								"account_classification": "Personal",
								"joint_account": false,
								"account_matching_opt_out": false,
								"secondary_identification": "A1B2C3D4"
							}
						}
		  			}`
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	created, err := c.Create(ctx, &Data{})

	// then
	require.Empty(t, err)
	require.Equal(t, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", created.Data.ID, "Invalid ID")
	require.Equal(t, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c", created.Data.OrganisationID, "Invalid OrganisationID")
	require.Equal(t, "accounts", created.Data.Type, "Invalid Type")
	require.Equal(t, int64(0), *created.Data.Version, "Invalid Version")
	require.Equal(t, "GB", created.Data.Attributes.Country, "Invalid Country")
	require.Equal(t, "GBP", created.Data.Attributes.BaseCurrency, "Invalid BaseCurrency")
	require.Equal(t, "41426819", created.Data.Attributes.AccountNumber, "Invalid AccountNumber")
	require.Equal(t, "400300", created.Data.Attributes.BankID, "Invalid BankID")
	require.Equal(t, "GBDSC", created.Data.Attributes.BankIDCode, "Invalid BankIDCode")
	require.Equal(t, "NWBKGB22", created.Data.Attributes.Bic, "Invalid Bic")
	require.Equal(t, "GB11NWBK40030041426819", created.Data.Attributes.Iban, "Invalid Iban")
	require.Equal(t, []string{"Samantha Holder"}, created.Data.Attributes.Name, "Invalid Name")
	require.Equal(t, []string{"Sam Holder"}, created.Data.Attributes.AlternativeNames, "Invalid AlternativeNames")
	require.Equal(t, "confirmed", *created.Data.Attributes.Status, "Invalid Status")
	require.Equal(t, []UserDefinedData{{Key: "Some account related key", Value: "Some account related value"}}, created.Data.Attributes.UserDefinedData, "Invalid UserDefinedData")
	require.Equal(t, "card", *created.Data.Attributes.ValidationType, "Invalid ValidationType")
	require.Equal(t, "############", *created.Data.Attributes.ReferenceMask, "Invalid ReferenceMask")
	require.Equal(t, "same_day", *created.Data.Attributes.AcceptanceQualifier, "Invalid AcceptanceQualifier")
	require.Equal(t, "Personal", *created.Data.Attributes.AccountClassification, "Invalid AccountClassification")
	require.Equal(t, false, *created.Data.Attributes.JointAccount, "Invalid JointAccount")
	require.Equal(t, false, *created.Data.Attributes.AccountMatchingOptOut, "Invalid AccountMatchingOptOut")
	require.Equal(t, "A1B2C3D4", *created.Data.Attributes.SecondaryIdentification, "Invalid SecondaryIdentification")
}

func Test_Accounts_CreateFailed_InvalidData(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 400
	responseBody := `{
						"error_message": "some error message", 
						"error_code":"b5930880-1001-453b-86bd-2c5e29bd98d7"
					}`
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	_, err := c.Create(ctx, NewWithGenID("orgID", &Attributes{}))

	// then
	require.IsType(t, &InvalidDataError{}, err, "Invalid error type")
	require.Equal(t, "some error message", err.(*InvalidDataError).Msg)
	require.Equal(t, "b5930880-1001-453b-86bd-2c5e29bd98d7", err.(*InvalidDataError).Code)
}

func Test_Accounts_CreateFailed_AlreadyExists(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 409
	responseBody := ``
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	_, err := c.Create(ctx, &Data{ID: "eb89cce1-3b1f-4b37-967f-23354c5ad61e"})

	// then
	require.IsType(t, &AccountAlreadyExistsError{}, err, "Invalid error type")
	require.Equal(t, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", err.(*AccountAlreadyExistsError).ID)
}

func Test_Accounts_FetchSuccess(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 200
	responseBody := `{
						"data": {
							"type": "accounts",
							"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
							"version": 0,
							"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
							"attributes": {
								"country": "GB"
							}
						}
		  			}`
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	fetched, err := c.Fetch(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")

	// then
	require.Empty(t, err)
	require.Equal(t, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", fetched.Data.ID, "Invalid ID")
	require.Equal(t, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c", fetched.Data.OrganisationID, "Invalid OrganisationID")
	require.Equal(t, "accounts", fetched.Data.Type, "Invalid Type")
	require.Equal(t, int64(0), *fetched.Data.Version, "Invalid Version")
	require.Equal(t, "GB", fetched.Data.Attributes.Country, "Invalid Country")
}

func Test_Accounts_FetchFailed_InvalidData(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 400
	responseBody := `{
						"error_message": "some error message", 
						"error_code":"b5930880-1001-453b-86bd-2c5e29bd98d7"
					}`
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	_, err := c.Fetch(ctx, "1")

	// then
	require.IsType(t, &InvalidDataError{}, err, "Invalid error type")
	require.Equal(t, "some error message", err.(*InvalidDataError).Msg)
	require.Equal(t, "b5930880-1001-453b-86bd-2c5e29bd98d7", err.(*InvalidDataError).Code)
}

func Test_Accounts_FetchFailed_UnknownAccount(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 404
	responseBody := ``
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	_, err := c.Fetch(ctx, "eb89cce1-3b1f-4b37-967f-23354c5ad61e")

	// then
	require.IsType(t, &AccountNotFoundError{}, err, "Invalid error type")
	require.Equal(t, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", err.(*AccountNotFoundError).ID)
}

func Test_Accounts_DeleteSuccess(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 204
	responseBody := ``
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	err := c.Delete(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 0)

	// then
	require.Empty(t, err)
}

func Test_Accounts_DeleteFailed_UnknownAccount(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 404
	responseBody := ``
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	err := c.Delete(ctx, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", 0)

	// then
	require.IsType(t, &AccountNotFoundError{}, err, "Invalid error type")
	require.Equal(t, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", err.(*AccountNotFoundError).ID)
}

func Test_Accounts_DeleteFailed_InvalidVersion(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 409
	responseBody := ``
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	err := c.Delete(ctx, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", 1)

	// then
	require.IsType(t, &InvalidVersionError{}, err, "Invalid error type")
	require.Equal(t, int64(1), err.(*InvalidVersionError).Ver)
}

func Test_Accounts_DeleteFailed_InvalidData(t *testing.T) {
	// given
	ctx := context.Background()
	responseCode := 400
	responseBody := `{
						"error_message": "some error message", 
						"error_code":"b5930880-1001-453b-86bd-2c5e29bd98d7"
					}`
	c := setUpMockClient(withResponse(responseCode, responseBody))

	// when
	err := c.Delete(ctx, "eb89cce1-3b1f-4b37-967f-23354c5ad61e", 1)

	// then
	require.IsType(t, &InvalidDataError{}, err, "Invalid error type")
	require.Equal(t, "some error message", err.(*InvalidDataError).Msg)
	require.Equal(t, "b5930880-1001-453b-86bd-2c5e29bd98d7", err.(*InvalidDataError).Code)
}

func setUpMockClient(r RoundTrip) Service {
	u, err := url.Parse("http://form3/v1")
	if err != nil {
		log.Fatal(err)
	}
	return NewClient(&http.Client{Transport: r}, *u)
}

// withResponse builds a RoundTrip function that returns HTTP response with given statusCode and body
func withResponse(statusCode int, body string) RoundTrip {
	return func(*http.Request) (*http.Response, error) {
		return buildResponse(statusCode, body), nil
	}
}

// buildResponse builds HTTP response with given statusCode and body
func buildResponse(statusCode int, respBody string) *http.Response {
	return &http.Response{
		StatusCode:    statusCode,
		Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
		ContentLength: int64(len(respBody)),
	}
}

type RoundTrip func(*http.Request) (*http.Response, error)

func (r RoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}
