package rest

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// content type
const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// paramPrep prepares parameters for the service.
func paramPrep(param string) string {
	return strings.TrimSpace(param)
}

func ReadPathParam(r *http.Request, param string) string {
	return paramPrep(chi.URLParam(r, param))
}

func ReadQueryParam(r *http.Request, param string) string {
	return paramPrep(r.URL.Query().Get(param))
}

func ReadHeader(r *http.Request, param string) string {
	return paramPrep(r.Header.Get(param))
}

// ReadBody reads JSON body object from a REST request.
func ReadBody(r *http.Request, object any) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return BadRequestErrorf("can't read body").WithError(err)
	}
	err = json.Unmarshal(body, object)
	if err != nil {
		return BadRequestErrorf("can't read body").WithError(err)
	}
	return nil
}

func WriteOK(w http.ResponseWriter, v any) error {
	return WriteJSON(w, v, http.StatusOK)
}

// SortingOrder represents the sorting order
type SortingOrder string

const (
	// SOAscending means ascending sorting order
	SOAscending SortingOrder = "asc"
	// SODescending means descending sorting order
	SODescending SortingOrder = "desc"
)

// Pagination of the lists.
type Pagination struct {
	// Pagination offset
	Offset int `json:"offset" validate:"min=0"`

	// Pagination limit
	Limit int `json:"limit" validate:"min=0,max=100"`

	// The field to use for sorting
	SortBy string `json:"sortBy"`

	// The order of sorting
	Order SortingOrder `json:"order" validate:"omitempty,oneof=asc desc"`
}

func ReadPaginationParams(r *http.Request) Pagination {
	limit, err := strconv.Atoi(ReadQueryParam(r, "limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(ReadQueryParam(r, "offset"))
	if err != nil {
		offset = 0
	}
	var order SortingOrder
	if ReadQueryParam(r, "order") != "desc" {
		order = SOAscending
	} else {
		order = SODescending
	}

	return Pagination{
		Limit:  limit,
		Offset: offset,
		SortBy: ReadQueryParam(r, "sortBy"),
		Order:  order,
	}
}
