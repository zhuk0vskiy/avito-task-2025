package controller

import (
	"encoding/json"

	"net/http"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	StatusCodeOuter int
}

func (w *StatusResponseWriter) WriteHeader(code int) {
	w.StatusCodeOuter = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *StatusResponseWriter) StatusCode() int {
	return w.StatusCodeOuter
}

type ErrorResponseStruct struct {
	Error string `json:"errors"`
}

// type SuccessResponseStruct struct {
// 	Data   interface{} `json:"data,omitempty"`
// }

func ErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponseStruct{Error: err})
}

func SuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}
