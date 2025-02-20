// Package message_manager provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package message_manager

import (
	externalRef0 "testoapi/gens/models/common"
)

// Defines values for MessageDirection.
const (
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"
)

// Defines values for MessageProviderName.
const (
	MessageProviderNameMessagebird MessageProviderName = "messagebird"
	MessageProviderNameTelnyx      MessageProviderName = "telnyx"
	MessageProviderNameTwilio      MessageProviderName = "twilio"
)

// Defines values for MessageType.
const (
	MessageTypeSMS MessageType = "sms"
)

// Defines values for TargetStatus.
const (
	TargetStatusDLRTimeout TargetStatus = "dlr_timeout"
	TargetStatusDelivered  TargetStatus = "delivered"
	TargetStatusFailed     TargetStatus = "failed"
	TargetStatusGWTimeout  TargetStatus = "gw_timeout"
	TargetStatusQueued     TargetStatus = "queued"
	TargetStatusReceived   TargetStatus = "received"
	TargetStatusSent       TargetStatus = "sent"
)

// Message defines model for Message.
type Message struct {
	// CustomerId Unique identifier for the customer.
	CustomerId *string `json:"customer_id,omitempty"`

	// Direction Direction of the message.
	Direction *MessageDirection `json:"direction,omitempty"`

	// Id Unique identifier for the message.
	Id *string `json:"id,omitempty"`

	// Source Contains source or destination detail info.
	Source *externalRef0.Address `json:"source,omitempty"`

	// Targets List of target addresses to which the message is sent.
	Targets *[]Target `json:"targets,omitempty"`

	// Text The text delivered in the body of the message.
	Text *string `json:"text,omitempty"`

	// TmCreate Timestamp when the message was created.
	TmCreate *string `json:"tm_create,omitempty"`

	// TmDelete Timestamp when the message was deleted.
	TmDelete *string `json:"tm_delete,omitempty"`

	// TmUpdate Timestamp when the message was last updated.
	TmUpdate *string `json:"tm_update,omitempty"`

	// Type Type of the message.
	Type *MessageType `json:"type,omitempty"`
}

// MessageDirection Direction of the message.
type MessageDirection string

// MessageProviderName Name of the message provider.
type MessageProviderName string

// MessageType Type of the message.
type MessageType string

// Target defines model for Target.
type Target struct {
	// Destination Contains source or destination detail info.
	Destination *externalRef0.Address `json:"destination,omitempty"`

	// Parts The number of message parts (if the message is split).
	Parts *int `json:"parts,omitempty"`

	// Status The status of the message for the target.
	Status *TargetStatus `json:"status,omitempty"`

	// TmUpdate Timestamp when the target message was last updated.
	TmUpdate *string `json:"tm_update,omitempty"`
}

// TargetStatus The status of the message for the target.
type TargetStatus string
