package internal

import (
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestStatus_API(t *testing.T) {
	w, r := rest.NewTestRequest()
	err := Status("0.0.0")(w, r)
	require.NoError(t, err)
	expected, err := json.Marshal(statusResponse{Status: "ok", Version: "0.0.0"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expected), w.Body.String())
}
