package post

type UpdatePostRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
