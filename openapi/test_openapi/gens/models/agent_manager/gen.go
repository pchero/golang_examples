// Package agent_manager provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package agent_manager

import (
	"time"

	externalRef0 "testoapi/gens/models/common"
)

// Defines values for AgentPermission.
const (
	PermissionAll               AgentPermission = 65535
	PermissionCustomerAdmin     AgentPermission = 32
	PermissionCustomerAgent     AgentPermission = 16
	PermissionCustomerAll       AgentPermission = 240
	PermissionCustomerManager   AgentPermission = 64
	PermissionNone              AgentPermission = 0
	PermissionProjectAll        AgentPermission = 15
	PermissionProjectSuperAdmin AgentPermission = 1
)

// Defines values for AgentRingMethod.
const (
	AgentRingMethodLinear  AgentRingMethod = "linear"
	AgentRingMethodRingAll AgentRingMethod = "ringall"
)

// Defines values for AgentStatus.
const (
	AgentStatusAvailable AgentStatus = "available"
	AgentStatusAway      AgentStatus = "away"
	AgentStatusBusy      AgentStatus = "busy"
	AgentStatusNone      AgentStatus = ""
	AgentStatusOffline   AgentStatus = "offline"
	AgentStatusRinging   AgentStatus = "ringing"
)

// Agent Represents an agent resource.
type Agent struct {
	// Addresses Agent's endpoint addresses.
	Addresses *[]externalRef0.Address `json:"addresses,omitempty"`

	// CustomerId Resource's customer ID.
	CustomerId *string `json:"customer_id,omitempty"`

	// Detail Agent's detail.
	Detail *string `json:"detail,omitempty"`
	Id     *string `json:"id,omitempty"`

	// Name Agent's name.
	Name *string `json:"name,omitempty"`

	// Permission Permission type
	Permission *AgentPermission `json:"permission,omitempty"`

	// RingMethod Represents an agent resource.
	RingMethod *AgentRingMethod `json:"ring_method,omitempty"`

	// Status Agent's status
	Status *AgentStatus `json:"status,omitempty"`

	// TagIds Agent's tag IDs.
	TagIds *[]string `json:"tag_ids,omitempty"`

	// TmCreate Created timestamp.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete Deleted timestamp.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmUpdate Updated timestamp.
	TmUpdate *time.Time `json:"tm_update,omitempty"`

	// Username Agent's username.
	Username *string `json:"username,omitempty"`
}

// AgentPermission Permission type
type AgentPermission uint64

// AgentRingMethod Represents an agent resource.
type AgentRingMethod string

// AgentStatus Agent's status
type AgentStatus string
