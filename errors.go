package golactus

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func NewError(code int, err string) *ResponseError {
	return &ResponseError{
		StatusCode: code,
		Message:    err,
	}
}

func (e *ResponseError) Error() string {
	return e.Message
}

func handleError(err error, w http.ResponseWriter) {
	respErr := &ResponseError{
		Message: err.Error(),
	}

	switch e := err.(type) {
	case *ResponseError:
		respErr = e
	default:
		respErr.StatusCode = http.StatusInternalServerError
	}

	jsonError(w, respErr)
}

func jsonError(w http.ResponseWriter, err *ResponseError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(err)
}
