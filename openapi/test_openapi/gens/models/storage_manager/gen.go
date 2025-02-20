// Package storage_manager provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package storage_manager

import (
	"time"
)

// Defines values for FilereferenceType.
const (
	Empty     FilereferenceType = ""
	Normal    FilereferenceType = "normal"
	Recording FilereferenceType = "recording"
)

// Account defines model for Account.
type Account struct {
	// CustomerId The customer ID associated with the account.
	CustomerId *string `json:"customer_id,omitempty"`

	// Id The unique identifier for the account.
	Id *string `json:"id,omitempty"`

	// TmCreate Timestamp when the account was created.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete Timestamp when the account was deleted.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmUpdate Timestamp when the account was last updated.
	TmUpdate *time.Time `json:"tm_update,omitempty"`

	// TotalFileCount The total number of files in the account.
	TotalFileCount *int64 `json:"total_file_count,omitempty"`

	// TotalFileSize The total file size in bytes.
	TotalFileSize *int64 `json:"total_file_size,omitempty"`
}

// File defines model for File.
type File struct {
	// CustomerId The customer ID associated with the file.
	CustomerId *string `json:"customer_id,omitempty"`

	// Detail The details of the file.
	Detail *string `json:"detail,omitempty"`

	// Filename The filename of the file.
	Filename *string `json:"filename,omitempty"`

	// Filesize The size of the file in bytes.
	Filesize *int64 `json:"filesize,omitempty"`

	// Id The unique identifier for the file.
	Id *string `json:"id,omitempty"`

	// Name The name of the file.
	Name *string `json:"name,omitempty"`

	// OwnerId The owner ID of the file.
	OwnerId *string `json:"owner_id,omitempty"`

	// ReferenceId The reference ID associated with the file.
	ReferenceId *string `json:"reference_id,omitempty"`

	// ReferenceType The reference type of the file.
	ReferenceType *FilereferenceType `json:"reference_type,omitempty"`

	// TmCreate Timestamp when the file was created.
	TmCreate *time.Time `json:"tm_create,omitempty"`

	// TmDelete Timestamp when the file was deleted.
	TmDelete *time.Time `json:"tm_delete,omitempty"`

	// TmDownloadExpire The timestamp when the download link expires.
	TmDownloadExpire *time.Time `json:"tm_download_expire,omitempty"`

	// TmUpdate Timestamp when the file was last updated.
	TmUpdate *time.Time `json:"tm_update,omitempty"`

	// UriDownload The URI for downloading the file.
	UriDownload *string `json:"uri_download,omitempty"`
}

// FilereferenceType The reference type of the file.
type FilereferenceType string
