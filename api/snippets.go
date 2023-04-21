package api

type createSnippetRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Expires int    `json:"expires"`
}
