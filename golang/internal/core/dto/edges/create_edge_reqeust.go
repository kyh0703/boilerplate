package edges

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

type CreateEdgeRequest struct {
	ID        uuid.UUID      `json:"id" validate:"required"`
	FlowID    int64          `json:"flow_id" validate:"required"`
	Source    int64          `json:"source" validate:"required"`
	Target    int64          `json:"target" validate:"required"`
	Type      string         `json:"type" validate:"required"`
	Label     sql.NullString `json:"label"`
	Hidden    sql.NullInt64  `json:"hidden"`
	MarkerEnd sql.NullString `json:"marker_end"`
	Points    sql.NullString `json:"points"`
}
