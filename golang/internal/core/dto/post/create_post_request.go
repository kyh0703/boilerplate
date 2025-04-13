package post

type CreatePostRequest struct {
	Title   string `json:"name" validate:"required"`
	Content string `json:"description" validate:"required"`
}
