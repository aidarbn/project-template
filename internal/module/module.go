package module

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	"bitbucket.org/creativeadvtech/project-template/internal/utils"
	"bitbucket.org/creativeadvtech/project-template/pkg/database"
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

type Module struct {
	db *bun.DB
}

func NewModule(db *bun.DB) *Module {
	return &Module{db: db}
}

const getListQuery = `
			SELECT *
			FROM (SELECT coalesce(json_agg(d.*), '[]'::json) as list
				  FROM (SELECT *
						FROM objects
						ORDER BY %s %s
						OFFSET %d LIMIT %d) as d) as data,
				 (SELECT count(*) as total
				  FROM objects as d) as total;
`

func (m Module) List(ctx context.Context, filter ListFilter) (*internal.ObjectList, error) {
	var list internal.ObjectList
	query := fmt.Sprintf(getListQuery, filter.SortBy, filter.Order, filter.Offset, filter.Limit)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	err = m.db.ScanRows(ctx, rows, &list)
	if err != nil {
		return nil, err
	}
	list.Count = len(list.List)
	return &list, nil
}

func (m Module) Get(ctx context.Context, id internal.UUID) (*internal.Object, error) {
	object := &internal.Object{ID: id}
	err := m.db.NewSelect().
		Model(object).
		WherePK().
		Scan(ctx)
	if database.IsNotFound(err) {
		return nil, rest.NotFoundErrorf(utils.TypeName(object) + " not found")
	} else if err != nil {
		return nil, err
	}
	return object, nil
}

func (m Module) Create(ctx context.Context, object *internal.Object) (*internal.Object, error) {
	object.ID = ""
	object.CreatedAt = time.Time{}
	object.UpdatedAt = time.Time{}
	res, err := m.db.NewInsert().
		Model(object).
		Returning("*").
		Exec(ctx)
	if database.IsDuplicate(err) {
		return nil, rest.BadRequestErrorf("duplicate " + utils.TypeName(object)).WithError(err)
	} else if err != nil {
		return nil, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return nil, err
	} else if rows != 1 {
		return nil, fmt.Errorf("only one row should be affected")
	}
	return object, nil
}

func (m Module) Update(ctx context.Context, object *internal.Object) (*internal.Object, error) {
	object.CreatedAt = time.Time{}
	object.UpdatedAt = time.Time{}
	res, err := m.db.NewUpdate().
		Model(object).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	} else if rows, err := res.RowsAffected(); err != nil {
		return nil, err
	} else if rows == 0 {
		return nil, rest.NotFoundErrorf(utils.TypeName(object) + " not found")
	} else if rows != 1 {
		return nil, fmt.Errorf("only one row should be affected")
	}
	return object, nil
}

func (m Module) Delete(ctx context.Context, id internal.UUID) error {
	object := &internal.Object{ID: id}
	res, err := m.db.NewDelete().
		Model(object).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	} else if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return rest.NotFoundErrorf(utils.TypeName(object) + " not found")
	} else if rows != 1 {
		return fmt.Errorf("only one row should be affected")
	}
	return nil
}
