package module

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	rest2 "bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"context"
	"net/http"
)

//go:generate mockery --name "Service" --inpackage --structname "mockService" --filename "service.mock.go"

type Service interface {
	List(ctx context.Context, filter ListFilter) (*internal.ObjectList, error)
	Get(ctx context.Context, id internal.UUID) (*internal.Object, error)
	Create(ctx context.Context, object *internal.Object) (*internal.Object, error)
	Update(ctx context.Context, object *internal.Object) (*internal.Object, error)
	Delete(ctx context.Context, id internal.UUID) error
}

type Rest struct {
	*rest2.Mux
	svc Service
}

func NewRest(svc Service) *Rest {
	res := &Rest{
		Mux: rest2.NewMux(),
		svc: svc,
	}

	res.Get("/", rest2.APIHandlerFunc(res.list))
	res.Get("/{ObjectID}", rest2.APIHandlerFunc(res.get))
	res.Post("/", rest2.APIHandlerFunc(res.create))
	res.Put("/{ObjectID}", rest2.APIHandlerFunc(res.update))
	res.Delete("/{ObjectID}", rest2.APIHandlerFunc(res.delete))

	return res
}

type ListFilter struct {
	rest2.Pagination `json:"inline"`
}

func (api Rest) list(w http.ResponseWriter, r *http.Request) error {
	filter := ListFilter{Pagination: rest2.ReadPaginationParams(r)}
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
	return rest2.WriteOK(w, objects)
}

func (api Rest) get(w http.ResponseWriter, r *http.Request) error {
	var id internal.UUID
	if err := internal.ParseUUID(rest2.ReadPathParam(r, "ObjectID"), &id); err != nil {
		return rest2.NotFoundErrorf("not found").WithError(err)
	}

	object, err := api.svc.Get(r.Context(), id)
	if err != nil {
		return err
	}
	return rest2.WriteOK(w, object)
}

func (api Rest) create(w http.ResponseWriter, r *http.Request) error {
	object := &internal.Object{}
	if err := rest2.ReadBody(r, object); err != nil {
		return rest2.BadRequestErrorf("can't parse body").WithError(err)
	}

	if err := api.PrepareParams(r.Context(), object); err != nil {
		return err
	}

	object, err := api.svc.Create(r.Context(), object)
	if err != nil {
		return err
	}
	return rest2.WriteOK(w, object)
}

func (api Rest) update(w http.ResponseWriter, r *http.Request) error {
	object := &internal.Object{}

	if err := rest2.ReadBody(r, object); err != nil {
		return rest2.BadRequestErrorf("can't parse body").WithError(err)
	}

	if err := internal.ParseUUID(rest2.ReadPathParam(r, "ObjectID"), &object.ID); err != nil {
		return rest2.NotFoundErrorf("not found").WithError(err)
	}

	if err := api.PrepareParams(r.Context(), object); err != nil {
		return err
	}

	object, err := api.svc.Update(r.Context(), object)
	if err != nil {
		return err
	}
	return rest2.WriteOK(w, object)
}

func (api Rest) delete(w http.ResponseWriter, r *http.Request) error {
	var id internal.UUID
	if err := internal.ParseUUID(rest2.ReadPathParam(r, "ObjectID"), &id); err != nil {
		return rest2.NotFoundErrorf("not found").WithError(err)
	}

	err := api.svc.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	return rest2.WriteOK(w, rest2.NewHTTPError(http.StatusOK, "successfully deleted"))
}
