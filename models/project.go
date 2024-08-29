package models

type Project struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description *string  `json:"description,omitempty"`
	Image       *string  `json:"image,omitempty"`
	Tag         []string `json:"tag,omitempty"`
	GitURL      *string  `json:"gitUrl,omitempty"`
	PreviewURL  *string  `json:"previewUrl,omitempty"`
}
