// Package todo provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package openapi

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// CreateList defines model for CreateList.
type CreateList struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	UserUUID    string `json:"userUUID"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// CreateListJSONBody defines parameters for CreateList.
type CreateListJSONBody CreateList

// CreateListJSONRequestBody defines body for CreateList for application/json ContentType.
type CreateListJSONRequestBody CreateListJSONBody
