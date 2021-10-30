package accounts

import "context"

// Accounts service interface
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about operations and fields.
type Service interface {

	// Register an existing bank account with Form3 or create a new one.
	// The country attribute must be specified as a minimum.
	// Depending on the country, other attributes such as bank_id and bic are mandatory.
	Create(ctx context.Context, account *Data) (*CreateSuccess, error)

	// Get a single account using the account ID.
	Fetch(ctx context.Context, id string) (*FetchSuccess, error)

	// Delete an account
	Delete(ctx context.Context, id string, version string) error
}

type CreateSuccess struct {
	Data  *Data  `json:"data"`
	Links *Links `json:"links"`
}

type FetchSuccess struct {
	Data  *Data  `json:"data"`
	Links *Links `json:"links"`
}

// Represents an account in the form3 org section.
type Data struct {
	// The specific attributes for each type of resource
	// Required: true
	Attributes *Attributes `json:"attributes,omitempty"`

	// The unique ID of the resource in UUID 4 format. It identifies the resource within the system.
	// Required: true
	ID string `json:"id,omitempty"`

	// The organisation ID of the organisation by which this resource has been created
	// Required: true
	OrganisationID string `json:"organisation_id,omitempty"`

	// The type of resource.
	Type string `json:"type,omitempty"`

	// A counter indicating how many times this resource has been modified.
	// Minimum: 0
	Version *int64 `json:"version,omitempty"`
}

type Attributes struct {
	// Classification of account, only used for Confirmation of Payee (CoP).
	// CoP: Can be either Personal or Business. Defaults to Personal.
	AccountClassification *string `json:"account_classification,omitempty"`

	//Flag to indicate if the account has opted out of account matching, only used for Confirmation of Payee
	// CoP: Set to true if the account has opted out of account matching. Defaults to false.
	AccountMatchingOptOut *bool `json:"account_matching_opt_out,omitempty"`

	// Account number. A unique account number will automatically be generated if not provided. If provided, the account number is not validated.
	AccountNumber string `json:"account_number,omitempty"`

	// Alternative primary account names, only used for UK Confirmation of Payee
	AlternativeNames []string `json:"alternative_names,omitempty"`

	// Local country bank identifier. Format depends on the country. Required for most countries.
	BankID string `json:"bank_id,omitempty"`

	// Identifies the type of bank ID being used, see documentation for allowed value for each country. Required value depends on country attribute.
	BankIDCode string `json:"bank_id_code,omitempty"`

	// ISO 4217 code used to identify the base currency of the account, e.g. 'GBP', 'EUR'
	BaseCurrency string `json:"base_currency,omitempty"`

	// SWIFT BIC in either 8 or 11 character format e.g. 'NWBKGB22'
	Bic string `json:"bic,omitempty"`

	// ISO 3166-1 code used to identify the domicile of the account, e.g. 'GB', 'FR'
	Country *string `json:"country,omitempty"`

	// A free-format reference that can be used to link this account to an external system
	CustomerID string `json:"customer_id,omitempty"`

	// IBAN of the account. Will be calculated from other fields if not supplied.
	Iban string `json:"iban,omitempty"`

	// Flag to indicate if the account is a joint account, only used for Confirmation of Payee (CoP).
	// CoP: Set to true is this is a joint account. Defaults to false.
	JointAccount *bool `json:"joint_account,omitempty"`

	// Name of the account holder, up to four lines possible.
	Name []string `json:"name,omitempty"`

	// Additional information to identify the account and account holder, only used for Confirmation of Payee (CoP)
	// CoP: Can be any type of additional identification, e.g. a building society roll number
	SecondaryIdentification string `json:"secondary_identification,omitempty"`

	// Status of the account
	// FPS: Can be pending, confirmed or closed. When this field is closed, status_reason must be provided.
	// SEPA & FPS Indirect (LHV): Can be either pending, confirmed or failed.
	// All other services: Can be pending or confirmed. pending is a virtual state and is immediately superseded by confirmed.
	Status *string `json:"status,omitempty"`

	// Provides additional account status information.
	// FPS: Must be a valid reason code if status is closed, cannot be provided otherwise.
	// Not supported for FPS Indirect via LHV.
	StatusReason *string `json:"status_reason,omitempty"`

	// Flag to indicate if the account has been switched away from this organisation, only used for Confirmation of Payee (CoP)
	// CoP: Set to true if the account has been switched using the Current Account Switching Service (CASS), false otherwise. Defaults to false.
	Switched *bool `json:"switched,omitempty"`
}

type Links struct {
	// Link to this endpoint or resource. Required.
	Self *string `json:"self"`
}

type ListLinks struct {
	Links

	// Link to the first page of a paginated response.
	First *string `json:"first,omitempty"`

	// Link to the last page of a paginated response.
	Last *string `json:"last,omitempty"`

	// Link to the next page of a paginated response.
	Next *string `json:"next,omitempty"`

	// Link to the previous page of a paginated response.
	Prev *string `json:"prev,omitempty"`
}
