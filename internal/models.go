package internal

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// UUID is an 32 hexadecimal digits that represents an ID
// swagger:strfmt uuid
type UUID string

func ParseUUID(s string, ID *UUID) error {
	if _, err := uuid.Parse(s); err != nil {
		return err
	}
	*ID = UUID(s)
	return nil
}

// Object is an object model for the module.
// swagger:model Object
type Object struct {
	bun.BaseModel `bun:"table:objects,alias:o"`
	// ID is a unique identifier of a user.
	// read only: true
	// unique: true
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID UUID `json:"id,omitempty" bun:"id,type:uuid,pk,nullzero,default:gen_random_uuid()" validate:"omitempty,uuid,required" mod:"trim"`

	// example: example data
	Data string `json:"data,omitempty" bun:"data" mod:"trim"`

	// read only: true
	CreatedAt time.Time `json:"created_at,omitempty" bun:"created_at,nullzero,notnull,default:now()"`

	// read only: true
	UpdatedAt time.Time `json:"updated_at,omitempty" bun:"updated_at,nullzero,notnull,default:now()"`
}

// ObjectList list of objects.
type ObjectList struct {
	// list of recorders
	// Read Only: true
	List []Object `json:"objects,omitempty"`

	// Number of items in a list
	// Read Only: true
	// Example: 1
	Count int `json:"count,omitempty"`

	// Total number of items available
	// Read Only: true
	// Example: 1
	Total int `json:"total,omitempty"`
}
