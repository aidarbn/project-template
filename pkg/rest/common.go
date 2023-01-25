package rest

import (
	"bitbucket.org/creativeadvtech/project-template/pkg/common"
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

func ReadPaginationParams(r *http.Request) common.Pagination {
	limit, err := strconv.Atoi(ReadQueryParam(r, "limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(ReadQueryParam(r, "offset"))
	if err != nil {
		offset = 0
	}
	var order common.SortingOrder
	if ReadQueryParam(r, "order") != "desc" {
		order = common.SOAscending
	} else {
		order = common.SODescending
	}

	return common.Pagination{
		Limit:  limit,
		Offset: offset,
		SortBy: ReadQueryParam(r, "sortBy"),
		Order:  order,
	}
}
