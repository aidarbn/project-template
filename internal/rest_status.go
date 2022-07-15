package internal

import (
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"net/http"
)

type statusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Status returns status of the service
func Status(version string) func(w http.ResponseWriter, _ *http.Request) error {
	return func(w http.ResponseWriter, _ *http.Request) error {
		return rest.WriteOK(w, statusResponse{Status: "ok", Version: version})
	}
}
