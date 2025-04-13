package entity

import "github.com/kyh0703/template/internal/pkg/constant"

type Group struct {
	ID     string `json:"id"`
	Extent bool   `json:"extent"`
	Expand string `json:"expand"`
}
type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Style struct {
	BgColor     string `json:"bgColor"`
	Color       string `json:"color"`
	BorderStyle string `json:"borderStyle"`
	BorderColor string `json:"borderColor"`
	Hidden      bool   `json:"hidden"`
}

type Node struct {
	ID       string            `json:"id"`
	FlowId   int64             `json:"flowId"`
	Type     constant.NodeType `json:"type"`
	Width    int32             `json:"width"`
	Height   int32             `json:"height"`
	Label    string            `json:"label"`
	Group    Group             `json:"group"`
	Position Position          `json:"position"`
	Style    Style             `json:"style"`
	Desc     string            `json:"desc"`
}
