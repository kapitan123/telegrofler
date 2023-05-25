// Package port provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package port

import (
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// Message defines model for Message.
type Message struct {
	Data        []byte             `json:"data"`
	Id          openapi_types.UUID `json:"id"`
	PublishTime time.Time          `json:"publishTime"`
}

// HandleVideoUrlPublishedMessageJSONRequestBody defines body for HandleVideoUrlPublishedMessage for application/json ContentType.
type HandleVideoUrlPublishedMessageJSONRequestBody = Message
