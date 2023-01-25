package object_module

import (
	"bitbucket.org/creativeadvtech/project-template/internal/models"
	"bitbucket.org/creativeadvtech/project-template/pkg/common"
	"bitbucket.org/creativeadvtech/project-template/pkg/errs"
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"context"
	"errors"
	"net/http"
)

//go:generate mockery --name "Service" --inpackage --structname "mockService" --filename "service.mock.go"

type Service interface {
	List(ctx context.Context, filter ListFilter) (*common.List[models.Object], error)
	Get(ctx context.Context, id common.UUID) (*models.Object, error)
	Create(ctx context.Context, object *createObject) (*models.Object, error)
	Update(ctx context.Context, id common.UUID, object *updateObject) (*models.Object, error)
	Delete(ctx context.Context, id common.UUID) error
}

type createObject struct {
	Data string `json:"data,omitempty" mod:"trim"`
}

type updateObject struct {
	Data string `json:"data,omitempty" mod:"trim"`
}

type Rest struct {
	*rest.Mux
	svc Service
}

func NewRest(svc Service) *Rest {
	res := &Rest{
		Mux: rest.NewMux(),
		svc: svc,
	}

	res.Get("/", rest.APIHandlerFunc(res.list))
	res.Get("/{ObjectID}", rest.APIHandlerFunc(res.get))
	res.Post("/", rest.APIHandlerFunc(res.create))
	res.Put("/{ObjectID}", rest.APIHandlerFunc(res.update))
	res.Delete("/{ObjectID}", rest.APIHandlerFunc(res.delete))

	return res
}

type ListFilter struct {
	common.Pagination `json:"inline"`
}

func (api Rest) list(w http.ResponseWriter, r *http.Request) error {
	filter := ListFilter{Pagination: rest.ReadPaginationParams(r)}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if err := api.PrepareParams(r.Context(), &filter); err != nil {
		return err
	}

	objects, err := api.svc.List(r.Context(), filter)
	if err != nil {
		return err
	}
	return rest.WriteOK(w, objects)
}

func (api Rest) get(w http.ResponseWriter, r *http.Request) error {
	var id common.UUID
	if err := common.ParseUUID(rest.ReadPathParam(r, "ObjectID"), &id); err != nil {
		return rest.NotFoundErrorf("not found").WithError(err)
	}

	object, err := api.svc.Get(r.Context(), id)
	var errNotFound *errs.NotFound
	if errors.As(err, &errNotFound) {
		return rest.NotFoundErrorf(errNotFound.Error())
	} else if err != nil {
		return err
	}
	return rest.WriteOK(w, object)
}

func (api Rest) create(w http.ResponseWriter, r *http.Request) error {
	cObject := &createObject{}
	if err := rest.ReadBody(r, cObject); err != nil {
		return rest.BadRequestErrorf("can't parse body").WithError(err)
	}

	if err := api.PrepareParams(r.Context(), cObject); err != nil {
		return err
	}

	object, err := api.svc.Create(r.Context(), cObject)
	var errDuplicate *errs.Duplicate
	if errors.As(err, &errDuplicate) {
		return rest.BadRequestErrorf(errDuplicate.Error())
	} else if err != nil {
		return err
	}
	return rest.WriteOK(w, object)
}

func (api Rest) update(w http.ResponseWriter, r *http.Request) error {
	uObject := &updateObject{}
	var id common.UUID

	if err := rest.ReadBody(r, uObject); err != nil {
		return rest.BadRequestErrorf("can't parse body").WithError(err)
	}

	if err := common.ParseUUID(rest.ReadPathParam(r, "ObjectID"), &id); err != nil {
		return rest.NotFoundErrorf("not found").WithError(err)
	}

	if err := api.PrepareParams(r.Context(), uObject); err != nil {
		return err
	}

	object, err := api.svc.Update(r.Context(), id, uObject)
	var errNotFound *errs.NotFound
	if errors.As(err, &errNotFound) {
		return rest.NotFoundErrorf(errNotFound.Error())
	} else if err != nil {
		return err
	}
	return rest.WriteOK(w, object)
}

func (api Rest) delete(w http.ResponseWriter, r *http.Request) error {
	var id common.UUID
	if err := common.ParseUUID(rest.ReadPathParam(r, "ObjectID"), &id); err != nil {
		return rest.NotFoundErrorf("not found").WithError(err)
	}

	err := api.svc.Delete(r.Context(), id)
	var errNotFound *errs.NotFound
	if errors.As(err, &errNotFound) {
		return rest.NotFoundErrorf(errNotFound.Error())
	} else if err != nil {
		return err
	}
	return rest.WriteOK(w, rest.NewHTTPError(http.StatusOK, "successfully deleted"))
}
