package models

type Note struct {
	ID      int    `json:"note_id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
