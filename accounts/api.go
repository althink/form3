package accounts

import (
	"context"

	"github.com/google/uuid"
)

// Accounts service interface
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about operations and fields.
type Service interface {

	// Create registers an existing bank account with Form3 or create a new one.
	// The country attribute must be specified as a minimum.
	// Depending on the country, other attributes such as bank_id and bic are mandatory.
	//
	// When data format is invalid returns InvalidDataError
	// When account with given id already exists returns AccountAlreadyExistsError
	// When other http error status returned returns HttpStatusError
	Create(ctx context.Context, account *Data) (*CreateSuccess, error)

	// Fetch returns a single account using the account ID.
	//
	// When accound with given id does not exist returns AccountNotFoundError
	// When data format is invalid returns InvalidDataError
	// When other http error status returned returns HttpStatusError
	Fetch(ctx context.Context, id string) (*FetchSuccess, error)

	// Delete deltes an account by given id and version.
	//
	// When accound with given id does not exist returns AccountNotFoundError
	// When version is invalid returns InvalidVersionError
	// When data format is invalid returns InvalidDataError
	// When other http error status returned returns HttpStatusError
	Delete(ctx context.Context, id string, version int64) error
}

type CreateSuccess struct {
	Data  *Data  `json:"data"`
	Links *Links `json:"links"`
}

type FetchSuccess struct {
	Data  *Data  `json:"data"`
	Links *Links `json:"links"`
}

// Create new account object
func New(id string, orgID string, attributes *Attributes) *Data {
	return &Data{
		ID:             id,
		OrganisationID: orgID,
		Type:           Type,
		Attributes:     attributes,
	}
}

// Create new account object with random ID
func NewWithGenID(orgID string, attributes *Attributes) *Data {
	return New(uuid.New().String(), orgID, attributes)
}

const Type = "accounts"

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

	// The type of resource: "accounts"
	Type string `json:"type,omitempty"`

	// A counter indicating how many times this resource has been modified.
	// Minimum: 0
	Version *int64 `json:"version,omitempty"`
}

type Attributes struct {

	// Determines the qualifier code with which payments to the account are accepted.
	// FPS: Must be a valid FPS Acceptance Qualifier if provided. Not supported for FPS Indirect via LHV.
	AcceptanceQualifier *string `json:"acceptance_qualifier,omitempty"`

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
	Country string `json:"country,omitempty"`

	// A free-format reference that can be used to link this account to an external system
	CustomerID string `json:"customer_id,omitempty"`

	// IBAN of the account. Will be calculated from other fields if not supplied.
	Iban string `json:"iban,omitempty"`

	// Flag to indicate if the account is a joint account, only used for Confirmation of Payee (CoP).
	// CoP: Set to true is this is a joint account. Defaults to false.
	JointAccount *bool `json:"joint_account,omitempty"`

	// Name of the account holder, up to four lines possible.
	Name []string `json:"name,omitempty"`

	// Mask to use when validating the reference field of inbound payments to the account
	// FPS: Payments without matching reference fields will be rejected.
	// ? matches any character, # matches any numeric characters (0-9), $ matches any alphanumeric character (a-z, A-Z). All other characters are literals.
	// \ can be used to escape control characters to literals. Maximum length 35 characters. - and space characters are ignored in the reference when matching.
	// Not supported for FPS Indirect via LHV.
	ReferenceMask *string `json:"reference_mask,omitempty"`

	// Additional information to identify the account and account holder, only used for Confirmation of Payee (CoP)
	// CoP: Can be any type of additional identification, e.g. a building society roll number
	SecondaryIdentification *string `json:"secondary_identification,omitempty"`

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

	// All-purpose list of key-value pairs to store data related to the account
	// FPS: Not supported for FPS Indirect via LHV.
	UserDefinedData []UserDefinedData `json:"user_defined_data,omitempty"`

	// Determines which validations are carried out on inbound payments to the account
	// FPS: Only card is allowed. Not supported for FPS Indirect via LHV.
	ValidationType *string `json:"validation_type,omitempty"`
}

type UserDefinedData struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
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
