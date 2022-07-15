package rest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructValidate(t *testing.T) {
	type TestStruct struct {
		Name     string `json:"name,omitempty" validate:"required"`
		LastName string `json:"last" validate:"required"`
	}
	t.Run("returns bad request error for the first invalid field", func(t *testing.T) {
		v := NewStructValidator()
		st := TestStruct{Name: ""}

		err := v.Validate(context.Background(), st)

		require.EqualError(t, err, "400: TestStruct.name: name is a required field")
	})
	t.Run("checks array elements too", func(t *testing.T) {
		type TestParentStruct struct {
			Elements []TestStruct `json:"elements" validate:"dive"`
		}
		v := NewStructValidator()
		st := TestParentStruct{
			Elements: []TestStruct{{}},
		}

		err := v.Validate(context.Background(), st)

		require.EqualError(t, err, "400: TestParentStruct.elements[0].name: name is a required field")
	})
	t.Run("ignores 'dive' if array is empty", func(t *testing.T) {
		type TestParentStruct struct {
			Elements []TestStruct `json:"elements" validate:"dive"`
		}
		v := NewStructValidator()
		st := TestParentStruct{
			Elements: []TestStruct{},
		}

		err := v.Validate(context.Background(), st)

		require.NoError(t, err)
	})
	t.Run("automatically checks memeber struct", func(t *testing.T) {
		type TestParentStruct struct {
			Member TestStruct `json:"member"`
		}
		v := NewStructValidator()
		st := TestParentStruct{
			Member: TestStruct{},
		}

		err := v.Validate(context.Background(), st)

		require.EqualError(t, err, "400: TestParentStruct.member.name: name is a required field")
	})
}
