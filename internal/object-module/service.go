package object_module

import (
	"bitbucket.org/creativeadvtech/project-template/internal/models"
	"bitbucket.org/creativeadvtech/project-template/pkg/common"
	"bitbucket.org/creativeadvtech/project-template/pkg/database"
	"bitbucket.org/creativeadvtech/project-template/pkg/errs"
	"context"
	"errors"
	"github.com/jinzhu/copier"
)

//go:generate mockery --name "Repository" --inpackage --structname "mockRepository" --filename "repository.mock.go"

type Repository interface {
	List(context.Context, common.Pagination) (*common.List[models.Object], error)
	Get(context.Context, common.UUID) (*models.Object, error)
	Create(context.Context, *models.Object) (*models.Object, error)
	Update(context.Context, common.UUID, *models.Object) (*models.Object, error)
	Delete(context.Context, common.UUID) error
}

type ObjectService struct {
	repo Repository
}

func NewModule(repo Repository) *ObjectService {
	return &ObjectService{repo: repo}
}

func (m ObjectService) List(ctx context.Context, filter ListFilter) (*common.List[models.Object], error) {
	return m.repo.List(ctx, filter.Pagination)
}

func (m ObjectService) Get(ctx context.Context, id common.UUID) (*models.Object, error) {
	resObj, err := m.repo.Get(ctx, id)
	if errors.Is(err, database.ErrNotFound) {
		return nil, errs.New[errs.NotFound]("object not found")
	}
	return resObj, err
}

func (m ObjectService) Create(ctx context.Context, object *createObject) (*models.Object, error) {
	var obj models.Object
	if err := copier.Copy(&obj, object); err != nil {
		return nil, err
	}

	resObj, err := m.repo.Create(ctx, &obj)
	if errors.Is(err, database.ErrDuplicate) {
		return nil, errs.New[errs.Duplicate]("object already exists")
	}
	return resObj, err
}

func (m ObjectService) Update(ctx context.Context, id common.UUID, object *updateObject) (*models.Object, error) {
	var obj models.Object
	if err := copier.Copy(&obj, object); err != nil {
		return nil, err
	}
	resObj, err := m.repo.Update(ctx, id, &obj)
	if errors.Is(err, database.ErrNotFound) {
		return nil, errs.New[errs.NotFound]("object not found")
	}
	return resObj, err
}

func (m ObjectService) Delete(ctx context.Context, id common.UUID) error {
	err := m.repo.Delete(ctx, id)
	if errors.Is(err, database.ErrNotFound) {
		return errs.New[errs.NotFound]("object not found")
	}
	return err
}
