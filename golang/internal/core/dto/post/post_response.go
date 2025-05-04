package post

type PostDto struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	UpdateAt string `json:"updateAt"`
	CreateAt string `json:"createAt"`
}
