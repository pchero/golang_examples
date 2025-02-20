// Package chatbot_manager provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package chatbot_manager

import (
	"time"
)

// Defines values for ChatbotEngineType.
const (
	ChatbotEngineTypeChatGPT ChatbotEngineType = "chatGPT"
	ChatbotEngineTypeClova   ChatbotEngineType = "clova"
)

// Defines values for ChatbotcallGender.
const (
	ChatbotcallGenderFemale  ChatbotcallGender = "female"
	ChatbotcallGenderMale    ChatbotcallGender = "male"
	ChatbotcallGenderNeutral ChatbotcallGender = "neutral"
)

// Defines values for ChatbotcallStatus.
const (
	ChatbotcallStatusEnd         ChatbotcallStatus = "end"
	ChatbotcallStatusInitiating  ChatbotcallStatus = "initiating"
	ChatbotcallStatusProgressing ChatbotcallStatus = "progressing"
)

// Defines values for ChatbotcallreferenceType.
const (
	ChatbotcallreferenceTypeCall ChatbotcallreferenceType = "call"
)

// Chatbot defines model for Chatbot.
type Chatbot struct {
	// CustomerId Unique identifier of the associated customer.
	CustomerId *string `json:"customer_id,omitempty"`

	// Detail Detailed information about the chatbot.
	Detail *string `json:"detail,omitempty"`

	// EngineType Type of engine used by the chatbot.
	EngineType *ChatbotEngineType `json:"engine_type,omitempty"`

	// Id Unique identifier of the chatbot.
	Id *string `json:"id,omitempty"`

	// InitPrompt Initial prompt to configure the chatbot's behavior.
	InitPrompt *string `json:"init_prompt,omitempty"`

	// Name Name of the chatbot.
	Name *string `json:"name,omitempty"`

	// TmCreate Timestamp when the chatbot was created.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete Timestamp when the chatbot was deleted.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmUpdate Timestamp when the chatbot was last updated.
	TmUpdate *time.Time `json:"tm_update,omitempty"`
}

// ChatbotEngineType Type of engine used by the chatbot.
type ChatbotEngineType string

// Chatbotcall defines model for Chatbotcall.
type Chatbotcall struct {
	// ActiveflowId Unique identifier for the active flow.
	ActiveflowId *string `json:"activeflow_id,omitempty"`

	// ChatbotId Unique identifier of the associated chatbot.
	ChatbotId *string `json:"chatbot_id,omitempty"`

	// ConfbridgeId Unique identifier for the conference bridge.
	ConfbridgeId *string `json:"confbridge_id,omitempty"`

	// CustomerId Unique identifier of the associated customer.
	CustomerId *string `json:"customer_id,omitempty"`

	// Gender Gender associated with the chatbot call.
	Gender *ChatbotcallGender `json:"gender,omitempty"`

	// Id Unique identifier for the chatbot call.
	Id *string `json:"id,omitempty"`

	// Language Language used during the chatbot call.
	Language *string `json:"language,omitempty"`

	// ReferenceId Unique identifier for the reference.
	ReferenceId *string `json:"reference_id,omitempty"`

	// ReferenceType Type of reference associated with the chatbot call.
	ReferenceType *ChatbotcallreferenceType `json:"reference_type,omitempty"`

	// Status Status of the chatbot call.
	Status *ChatbotcallStatus `json:"status,omitempty"`

	// TmCreate Timestamp when the chatbot call was created.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete Timestamp when the chatbot call was deleted.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmEnd Timestamp when the chatbot call ended.
	TmEnd *time.Time `json:"tm_end,omitempty"`

	// TmUpdate Timestamp when the chatbot call was last updated.
	TmUpdate *time.Time `json:"tm_update,omitempty"`

	// TranscribeId Unique identifier for the transcription service.
	TranscribeId *string `json:"transcribe_id,omitempty"`
}

// ChatbotcallGender Gender associated with the chatbot call.
type ChatbotcallGender string

// ChatbotcallMessage defines model for ChatbotcallMessage.
type ChatbotcallMessage struct {
	// Content Content of the message.
	Content *string `json:"content,omitempty"`

	// Role Role of the entity in the conversation.
	Role *string `json:"role,omitempty"`
}

// ChatbotcallStatus Status of the chatbot call.
type ChatbotcallStatus string

// ChatbotcallreferenceType Type of reference associated with the chatbot call.
type ChatbotcallreferenceType string
