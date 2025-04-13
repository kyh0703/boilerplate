package post

type PostResponse struct {
	Title    string `json:"name"`
	Content  string `json:"description"`
	UpdateAt string `json:"updateAt"`
	CreateAt string `json:"createAt"`
}
