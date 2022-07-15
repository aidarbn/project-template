package module

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	rest2 "bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

const (
	testID internal.UUID = "123e4567-e89b-12d3-a456-426655440000"
)

func testObject() *internal.Object {
	return &internal.Object{
		ID:        testID,
		Data:      "some data",
		CreatedAt: time.Date(2022, 07, 02, 00, 00, 00, 00, time.UTC),
		UpdatedAt: time.Date(2022, 07, 02, 00, 00, 00, 00, time.UTC),
	}
}

func testRecorderList() *internal.ObjectList {
	return &internal.ObjectList{
		List: []internal.Object{
			*testObject(),
		},
		Count: 1,
		Total: 10,
	}
}

func TestResource_list(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("List", mock.Anything,
			ListFilter{
				Pagination: rest2.Pagination{
					Offset: 0,
					Limit:  1,
					SortBy: "created_at",
					Order:  "asc",
				},
			}).
			Return(testRecorderList(), nil)

		w, r := rest2.NewTestRequest(
			rest2.WithQuery("limit", 1),
			rest2.WithQuery("offset", 0),
			rest2.WithQuery("sortBy", "created_at"),
			rest2.WithQuery("order", "asc"),
		)

		err := res.list(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(testRecorderList())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("List", mock.Anything,
			ListFilter{
				Pagination: rest2.Pagination{
					Offset: 0,
					Limit:  1,
					SortBy: "created_at",
					Order:  "asc",
				},
			}).
			Return(nil, fmt.Errorf("some error"))

		w, r := rest2.NewTestRequest(
			rest2.WithQuery("limit", 1),
			rest2.WithQuery("offset", 0),
			rest2.WithQuery("sortBy", "created_at"),
			rest2.WithQuery("order", "asc"),
		)

		err := res.list(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestResource_get(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Get", mock.Anything, testID).
			Return(testObject(), nil)

		w, r := rest2.NewTestRequest(
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.get(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(testObject())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Get", mock.Anything, testID).
			Return(nil, fmt.Errorf("some error"))

		w, r := rest2.NewTestRequest(
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.get(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestResource_create(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Create", mock.Anything, testObject()).
			Return(testObject(), nil)

		w, r := rest2.NewTestRequest(
			rest2.WithJSON(testObject()),
		)

		err := res.create(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(testObject())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Create", mock.Anything, testObject()).
			Return(nil, fmt.Errorf("some error"))

		w, r := rest2.NewTestRequest(
			rest2.WithJSON(testObject()),
		)

		err := res.create(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestResource_update(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Update", mock.Anything, testObject()).
			Return(testObject(), nil)

		w, r := rest2.NewTestRequest(
			rest2.WithJSON(testObject()),
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.update(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(testObject())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Update", mock.Anything, testObject()).
			Return(nil, fmt.Errorf("some error"))

		w, r := rest2.NewTestRequest(
			rest2.WithJSON(testObject()),
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.update(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestResource_delete(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Delete", mock.Anything, testID).
			Return(nil)

		w, r := rest2.NewTestRequest(
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.delete(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(map[string]any{"code": 200, "description": "successfully deleted"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Delete", mock.Anything, testID).
			Return(fmt.Errorf("some error"))

		w, r := rest2.NewTestRequest(
			rest2.WithPathParam("ObjectID", testID),
		)

		err := res.delete(w, r)
		require.EqualError(t, err, "some error")
	})
}
