package dto

import "time"

type GenerateImageRequest struct {
	Prompts []string `json:"prompts" binding:"required"`
	Width   int      `json:"width" binding:"required,min=256,max=2048"`
	Height  int      `json:"height" binding:"required,min=256,max=2048"`
	Quality string   `json:"quality" binding:"oneof=standard hd"`
	Style   string   `json:"style"`
}

type ImageResult struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Prompt    string    `json:"prompt"`
	CreatedAt time.Time `json:"created_at"`
}

type GenerateImageResponse struct {
	Results []ImageResult `json:"results"`
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
}

type ImageProgressResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	URL      string `json:"url,omitempty"`
	Error    string `json:"error,omitempty"`
}
