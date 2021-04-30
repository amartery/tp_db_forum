package models

import "github.com/go-openapi/strfmt"

//easyjson:json
type Thread struct {
	ID      int             `json:"id"`
	Forum   string          `json:"forum"`
	Title   string          `json:"title"`
	Author  string          `json:"author"`
	Message string          `json:"message"`
	Slug    *string         `json:"slug,omitempty"`
	Created strfmt.DateTime `json:"created,omitempty"`
	Votes   int             `json:"votes"`
}
