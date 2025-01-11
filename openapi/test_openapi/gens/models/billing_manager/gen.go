// Package billing_manager provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package billing_manager

import (
	"time"
)

// Defines values for AccountPaymentMethod.
const (
	AccountPaymentMethodCreditCard AccountPaymentMethod = "credit card"
	AccountPaymentMethodNone       AccountPaymentMethod = ""
)

// Defines values for AccountPaymentType.
const (
	AccountPaymentTypeNone    AccountPaymentType = ""
	AccountPaymentTypePrepaid AccountPaymentType = "prepaid"
)

// Defines values for BillingReferenceType.
const (
	BillingReferenceTypeCall        BillingReferenceType = "call"
	BillingReferenceTypeNone        BillingReferenceType = ""
	BillingReferenceTypeNumber      BillingReferenceType = "number"
	BillingReferenceTypeNumberRenew BillingReferenceType = "number_renew"
	BillingReferenceTypeSMS         BillingReferenceType = "sms"
)

// Defines values for BillingStatus.
const (
	BillingStatusEnd         BillingStatus = "end"
	BillingStatusFinished    BillingStatus = "finished"
	BillingStatusPending     BillingStatus = "pending"
	BillingStatusProgressing BillingStatus = "progressing"
)

// Account defines model for Account.
type Account struct {
	// Balance The balance of the account in USD.
	Balance *float32 `json:"balance,omitempty"`

	// CustomerId The unique identifier of the associated customer.
	CustomerId *string `json:"customer_id,omitempty"`

	// Detail Additional details about the account.
	Detail *string `json:"detail,omitempty"`

	// Id The unique identifier of the account.
	Id *string `json:"id,omitempty"`

	// Name The name of the account.
	Name *string `json:"name,omitempty"`

	// PaymentMethod The method of payment used for the account.
	PaymentMethod *AccountPaymentMethod `json:"payment_method,omitempty"`

	// PaymentType The type of payment associated with the account.
	PaymentType *AccountPaymentType `json:"payment_type,omitempty"`

	// TmCreate The timestamp when the account was created.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete The timestamp when the account was deleted, if applicable.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmUpdate The timestamp when the account was last updated.
	TmUpdate *time.Time `json:"tm_update,omitempty"`
}

// AccountPaymentMethod The method of payment used for the account.
type AccountPaymentMethod string

// AccountPaymentType The type of payment associated with the account.
type AccountPaymentType string

// Billing defines model for Billing.
type Billing struct {
	// AccountId The billing account ID.
	AccountId *string `json:"account_id,omitempty"`

	// BillingUnitCount The total count of billing units.
	BillingUnitCount *float32 `json:"billing_unit_count,omitempty"`

	// CostPerUnit The cost per billing unit.
	CostPerUnit *float32 `json:"cost_per_unit,omitempty"`

	// CostTotal The total cost of this billing.
	CostTotal *float32 `json:"cost_total,omitempty"`

	// CustomerId The customer's unique identifier.
	CustomerId *string `json:"customer_id,omitempty"`

	// Id The unique identifier of the billing.
	Id *string `json:"id,omitempty"`

	// ReferenceId The ID of the reference related to this billing.
	ReferenceId *string `json:"reference_id,omitempty"`

	// ReferenceType The type of reference associated with this billing.
	ReferenceType *BillingReferenceType `json:"reference_type,omitempty"`

	// Status Status of the billing.
	Status *BillingStatus `json:"status,omitempty"`

	// TmBillingEnd The end timestamp of the billing period.
	TmBillingEnd *time.Time `json:"tm_billing_end,omitempty"`

	// TmBillingStart The start timestamp of the billing period.
	TmBillingStart *time.Time `json:"tm_billing_start,omitempty"`

	// TmCreate The creation timestamp.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete The deletion timestamp, if applicable.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmUpdate The last update timestamp.
	TmUpdate *time.Time `json:"tm_update,omitempty"`
}

// BillingReferenceType The type of reference associated with this billing.
type BillingReferenceType string

// BillingStatus Status of the billing.
type BillingStatus string
