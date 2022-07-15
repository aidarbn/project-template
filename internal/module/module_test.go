package module

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	"bitbucket.org/creativeadvtech/project-template/internal/config"
	"bitbucket.org/creativeadvtech/project-template/pkg/database"
	"bitbucket.org/creativeadvtech/project-template/pkg/logging"
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/dbfixture"
	"os"
	"testing"
	"time"
)

func prepareFixture(t *testing.T) *Module {
	var cfg config.TestConfig
	err := envconfig.Process("", &cfg)
	require.NoError(t, err)
	db, err := database.NewDatabase(fmt.Sprintf(database.DbURLTemplate, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName), true)
	require.NoError(t, err)
	db.RegisterModel((*internal.Object)(nil))
	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	err = fixture.Load(context.Background(), os.DirFS("."), "module-dbfixture.yml")
	require.NoError(t, err)
	return NewModule(db)
}

func TestModule_List_Integration(t *testing.T) {
	if testing.Short() {
		return
	}
	logging.Init("debug")
	srv := prepareFixture(t)
	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		objects, err := srv.List(ctx,
			ListFilter{
				Pagination: rest.Pagination{
					Offset: 1,
					Limit:  3,
					SortBy: "id",
					Order:  "asc",
				},
			},
		)
		require.NoError(t, err)
		assert.Equal(t, 3, objects.Count)
		assert.Equal(t, 6, objects.Total)
		assert.Len(t, objects.List, 3)
		assert.Equal(t, internal.UUID("56d1c6bf-9d3d-4265-8be3-3a5b8b356291"), objects.List[0].ID)
		assert.Equal(t, "test data 1", objects.List[0].Data)
		assert.Equal(t, internal.UUID("56d1c6bf-9d3d-4265-8be3-3a5b8b356292"), objects.List[1].ID)
		assert.Equal(t, "test data 2", objects.List[1].Data)
		assert.Equal(t, internal.UUID("56d1c6bf-9d3d-4265-8be3-3a5b8b356293"), objects.List[2].ID)
		assert.Equal(t, "test data 3", objects.List[2].Data)
	})
}

func TestModule_Get_Integration(t *testing.T) {
	if testing.Short() {
		return
	}
	logging.Init("debug")
	srv := prepareFixture(t)
	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		obj, err := srv.Get(ctx, "56d1c6bf-9d3d-4265-8be3-3a5b8b356291")
		require.NoError(t, err)
		assert.Equal(t, internal.UUID("56d1c6bf-9d3d-4265-8be3-3a5b8b356291"), obj.ID)
		assert.Equal(t, "test data 1", obj.Data)
	})

	t.Run("error; not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := srv.Get(ctx, "56d1c6bf-9d3d-4265-8be3-3a5b8b356299")
		assert.EqualError(t, err, "404: Object not found")
	})
}

func TestModule_Create_Integration(t *testing.T) {
	if testing.Short() {
		return
	}
	logging.Init("debug")
	srv := prepareFixture(t)
	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		obj, err := srv.Create(ctx, &internal.Object{Data: "create test data"})
		require.NoError(t, err)
		assert.NotEmpty(t, obj.ID)
		assert.Equal(t, "create test data", obj.Data)
	})
}

func TestModule_Update_Integration(t *testing.T) {
	if testing.Short() {
		return
	}
	logging.Init("debug")
	srv := prepareFixture(t)
	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		obj, err := srv.Update(ctx, &internal.Object{ID: "56d1c6bf-9d3d-4265-8be3-3a5b8b356291", Data: "update test data"})
		require.NoError(t, err)
		assert.NotEmpty(t, obj.ID)
		assert.Equal(t, "update test data", obj.Data)
	})

	t.Run("error; not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := srv.Update(ctx, &internal.Object{ID: "56d1c6bf-9d3d-4265-8be3-3a5b8b356299", Data: "update test data"})
		assert.EqualError(t, err, "404: Object not found")
	})
}

func TestModule_Delete_Integration(t *testing.T) {
	if testing.Short() {
		return
	}
	logging.Init("debug")
	srv := prepareFixture(t)
	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Delete(ctx, "56d1c6bf-9d3d-4265-8be3-3a5b8b356291")
		require.NoError(t, err)
	})
	t.Run("error; not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Delete(ctx, "56d1c6bf-9d3d-4265-8be3-3a5b8b356299")
		assert.EqualError(t, err, "404: Object not found")
	})
}
