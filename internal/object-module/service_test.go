package object_module

import (
	"bitbucket.org/creativeadvtech/project-template/internal/models"
	"bitbucket.org/creativeadvtech/project-template/pkg/common"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestModule_List(t *testing.T) {
	repo := &mockRepository{}
	srv := NewModule(repo)

	t.Run("OK", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("List", mock.Anything,
			common.Pagination{
				Offset: 1,
				Limit:  3,
				SortBy: "id",
				Order:  "asc",
			}).
			Return(testObjectList(), nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		list, err := srv.List(ctx, ListFilter{
			Pagination: common.Pagination{
				Offset: 1,
				Limit:  3,
				SortBy: "id",
				Order:  "asc",
			},
		})

		require.NoError(t, err)
		assert.Equal(t, testObjectList(), list)
	})

	t.Run("Error", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("List", mock.Anything,
			common.Pagination{
				Offset: 1,
				Limit:  3,
				SortBy: "id",
				Order:  "asc",
			}).
			Return(nil, fmt.Errorf("some error"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := srv.List(ctx, ListFilter{
			Pagination: common.Pagination{
				Offset: 1,
				Limit:  3,
				SortBy: "id",
				Order:  "asc",
			},
		})

		require.EqualError(t, err, "some error")
	})
}

func TestModule_Get(t *testing.T) {
	repo := &mockRepository{}
	srv := NewModule(repo)

	t.Run("OK", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Get", mock.Anything, testID).
			Return(testObject(), nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		object, err := srv.Get(ctx, testID)

		require.NoError(t, err)
		assert.Equal(t, testObject(), object)
	})

	t.Run("Error", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Get", mock.Anything, testID).
			Return(nil, fmt.Errorf("some error"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := srv.Get(ctx, testID)

		require.EqualError(t, err, "some error")
	})
}

func TestModule_Create(t *testing.T) {
	repo := &mockRepository{}
	srv := NewModule(repo)

	t.Run("OK", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Create", mock.Anything, &models.Object{Data: "some data"}).
			Return(testObject(), nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		object, err := srv.Create(ctx, testCreateObject())

		require.NoError(t, err)
		assert.Equal(t, testObject(), object)
	})

	t.Run("Error", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Create", mock.Anything, &models.Object{Data: "some data"}).
			Return(nil, fmt.Errorf("some error"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := srv.Create(ctx, testCreateObject())

		require.EqualError(t, err, "some error")
	})
}

func TestModule_Update(t *testing.T) {
	repo := &mockRepository{}
	srv := NewModule(repo)

	t.Run("OK", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Update", mock.Anything, testID, &models.Object{Data: "some data"}).
			Return(testObject(), nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		object, err := srv.Update(ctx, testID, testUpdateObject())

		require.NoError(t, err)
		assert.Equal(t, testObject(), object)
	})

	t.Run("Error", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Update", mock.Anything, testID, &models.Object{Data: "some data"}).
			Return(nil, fmt.Errorf("some error"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := srv.Update(ctx, testID, testUpdateObject())

		require.EqualError(t, err, "some error")
	})
}

func TestModule_Delete(t *testing.T) {
	repo := &mockRepository{}
	srv := NewModule(repo)

	t.Run("OK", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Delete", mock.Anything, testID).
			Return(nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Delete(ctx, testID)

		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		repo.Mock = mock.Mock{}

		repo.On("Delete", mock.Anything, testID).
			Return(fmt.Errorf("some error"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Delete(ctx, testID)

		require.EqualError(t, err, "some error")
	})
}
