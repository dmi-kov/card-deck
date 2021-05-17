package httphelper

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteErrorResponse(t *testing.T) {
	var rw http.ResponseWriter

	err := WriteErrorResponse(rw, 200, "error msg")
	require.EqualError(t, err, "not valid status code 200 for error response")
}

func TestWriteSuccessResponse(t *testing.T) {
	var rw http.ResponseWriter

	err := WriteSuccessResponse(rw, 401, "payload")
	require.EqualError(t, err, "not valid status code 401 for success response")
}
