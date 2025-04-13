package flows

type FlowResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdateAt    string `json:"updateAt"`
	CreateAt    string `json:"createAt"`
}
