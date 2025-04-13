package entity

import (
	"encoding/json"

	"github.com/kyh0703/template/internal/core/domain/model"
)

type MarkerEnd struct{}

type Point struct{}

type Edge struct {
	ID        int64
	FlowID    int64
	Source    string
	Target    string
	Hidden    bool
	MarkerEnd *MarkerEnd
	Points    []Point
}

func (e *Edge) ToModel() (*model.Edge, error) {
	var model model.Edge
	model.ID = e.ID
	model.FlowID = e.FlowID
	model.Source = e.Source
	model.Target = e.Target
	if e.Hidden {
		model.Hidden = 1
	} else {
		model.Hidden = 0
	}
	if e.MarkerEnd != nil {
		markerEnd, err := json.Marshal(e.MarkerEnd)
		if err != nil {
			return nil, err
		}
		model.MarkerEnd = string(markerEnd)
	}
	return &model, nil
}

func (e *Edge) FromModel(model *model.Edge) {
	e.ID = model.ID
	e.FlowID = model.FlowID
	e.Source = model.Source
	e.Target = model.Target
	e.Hidden = model.Hidden == 1
	if model.MarkerEnd != "" {
		var markerEnd MarkerEnd
		err := json.Unmarshal([]byte(model.MarkerEnd), &markerEnd)
		if err == nil {
			e.MarkerEnd = &markerEnd
		} else {
			e.MarkerEnd = nil
		}
	}
}
