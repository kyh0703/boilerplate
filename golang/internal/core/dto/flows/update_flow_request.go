package flows

type UpdateFlowRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
