package responseModel

import "time"

type AuditFields struct {
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	CreatedById *uint      `json:"createdById,omitempty"`
	UpdatedById *uint      `json:"updatedById,omitempty"`
}
