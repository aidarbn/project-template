package object_module

import (
	"bitbucket.org/creativeadvtech/project-template/internal/models"
	"bitbucket.org/creativeadvtech/project-template/pkg/common"
	"bitbucket.org/creativeadvtech/project-template/pkg/testutils"
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
	testID common.UUID = "123e4567-e89b-12d3-a456-426655440000"
)

func testObject() *models.Object {
	return &models.Object{
		ID:        testID,
		Data:      "some data",
		CreatedAt: time.Date(2022, 07, 02, 00, 00, 00, 00, time.UTC),
		UpdatedAt: time.Date(2022, 07, 02, 00, 00, 00, 00, time.UTC),
	}
}

func testCreateObject() *createObject {
	return &createObject{
		Data: "some data",
	}
}

func testUpdateObject() *updateObject {
	return &updateObject{
		Data: "some data",
	}
}

func testObjectList() *common.List[models.Object] {
	return &common.List[models.Object]{
		List: []models.Object{
			*testObject(),
		},
		Count: 1,
		Total: 10,
	}
}

func TestRest_list(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("List", mock.Anything,
			ListFilter{
				Pagination: common.Pagination{
					Offset: 0,
					Limit:  1,
					SortBy: "created_at",
					Order:  "asc",
				},
			}).
			Return(testObjectList(), nil)

		w, r := testutils.NewTestRequest(
			testutils.WithQuery("limit", 1),
			testutils.WithQuery("offset", 0),
			testutils.WithQuery("sortBy", "created_at"),
			testutils.WithQuery("order", "asc"),
		)

		err := res.list(w, r)
		require.NoError(t, err)
		expected, err := json.Marshal(testObjectList())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expected), w.Body.String())
	})

	t.Run("Error", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("List", mock.Anything,
			ListFilter{
				Pagination: common.Pagination{
					Offset: 0,
					Limit:  1,
					SortBy: "created_at",
					Order:  "asc",
				},
			}).
			Return(nil, fmt.Errorf("some error"))

		w, r := testutils.NewTestRequest(
			testutils.WithQuery("limit", 1),
			testutils.WithQuery("offset", 0),
			testutils.WithQuery("sortBy", "created_at"),
			testutils.WithQuery("order", "asc"),
		)

		err := res.list(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestRest_get(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Get", mock.Anything, testID).
			Return(testObject(), nil)

		w, r := testutils.NewTestRequest(
			testutils.WithPathParam("ObjectID", testID),
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

		w, r := testutils.NewTestRequest(
			testutils.WithPathParam("ObjectID", testID),
		)

		err := res.get(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestRest_create(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Create", mock.Anything, testCreateObject()).
			Return(testObject(), nil)

		w, r := testutils.NewTestRequest(
			testutils.WithJSON(testObject()),
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

		srv.On("Create", mock.Anything, testCreateObject()).
			Return(nil, fmt.Errorf("some error"))

		w, r := testutils.NewTestRequest(
			testutils.WithJSON(testObject()),
		)

		err := res.create(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestRest_update(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Update", mock.Anything, testID, testUpdateObject()).
			Return(testObject(), nil)

		w, r := testutils.NewTestRequest(
			testutils.WithJSON(testObject()),
			testutils.WithPathParam("ObjectID", testID),
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

		srv.On("Update", mock.Anything, testID, testUpdateObject()).
			Return(nil, fmt.Errorf("some error"))

		w, r := testutils.NewTestRequest(
			testutils.WithJSON(testObject()),
			testutils.WithPathParam("ObjectID", testID),
		)

		err := res.update(w, r)
		require.EqualError(t, err, "some error")
	})
}

func TestRest_delete(t *testing.T) {
	srv := &mockService{}
	res := NewRest(srv)
	t.Run("OK", func(t *testing.T) {
		srv.Mock = mock.Mock{}

		srv.On("Delete", mock.Anything, testID).
			Return(nil)

		w, r := testutils.NewTestRequest(
			testutils.WithPathParam("ObjectID", testID),
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

		w, r := testutils.NewTestRequest(
			testutils.WithPathParam("ObjectID", testID),
		)

		err := res.delete(w, r)
		require.EqualError(t, err, "some error")
	})
}
