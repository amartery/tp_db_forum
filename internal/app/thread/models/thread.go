package models

//easyjson:json
type Thread struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created"`
}

type UpdateRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}