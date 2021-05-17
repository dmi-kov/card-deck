package httphelper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/card-deck/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Handler is a type that wraps any HTTP handler function
// to satisfy the http.Handler interface.
type Handler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP handles the HTTP request.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		var data interface{}
		var status int
		switch err := err.(type) {
		case *errors.Error:
			status = int(err.Kind())
			data = err.Context()
		default:
			status = http.StatusInternalServerError
			data = "Oops! Something went wrong"
		}
		zap.L().Error("error occurred while handle request", zap.String("error", err.Error()))
		if err := WriteErrorResponse(w, status, data); err != nil {
			logrus.Errorf("unable to write the json response: %v", err)
		}
	}
}

// WriteSuccessResponse writes a success json response
func WriteSuccessResponse(w http.ResponseWriter, status int, payload interface{}) error {
	if status >= 200 && status <= 299 {
		return WriteJSON(w, status, payload)
	}
	return errors.Newf(errors.Internal, "not valid status code %d for success response", status)
}

// WriteErrorResponse writes an error json response
func WriteErrorResponse(w http.ResponseWriter, status int, error interface{}) error {
	if status >= 400 && status <= 599 {
		resp := map[string]interface{}{
			"error": error,
		}
		return WriteJSON(w, status, resp)
	}
	return errors.Newf(errors.Internal, "not valid status code %d for error response", status)
}

// WriteJSON writes a json response
func WriteJSON(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return errors.Wrapf(err, errors.Internal, "failed to write the response body")
	}
	return nil
}

// ReadJSON reads a json request body.
func ReadJSON(r *http.Request, receiver interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.InvalidInput.Wrap(err, "the request body is missing")
	}
	if err := json.Unmarshal(body, receiver); err != nil {
		logrus.Error(err)
		return errors.InvalidInput.Wrap(err, "the request body is not valid json")
	}
	return nil
}
