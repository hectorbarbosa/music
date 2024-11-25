package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"music/internal"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	// Error string
	Error string `json:"error" example:"error description"`
}

func renderErrorResponse(w http.ResponseWriter, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var ierr *internal.Error
	if !errors.As(err, &ierr) {
		resp.Error = "internal error"
	} else {
		switch ierr.Code() {
		case internal.ErrorCodeNotFound:
			status = http.StatusNotFound
		case internal.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest
		case internal.ErrorCodeBadGateWay:
			status = http.StatusBadGateway
		}
	}

	renderResponse(w, resp, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		fmt.Println("error writing content")
	}
}
