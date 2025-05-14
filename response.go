package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type ResponseBody struct {
	StatusCode int    `json:"-"`
	Success    bool   `json:"success"`
	Data       any    `json:"data,omitempty"`
	Err        string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
	Timestamp  int64  `json:"ts"`
}

func (resBody *ResponseBody) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resBody.StatusCode)
	resBody.Timestamp = time.Now().Unix()

	return nil
}

func NewSuccessResponse(data any) *ResponseBody {
	return &ResponseBody{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
}

func NewErrorResponse(err error) *ResponseBody {
	var restErr *Error
	if !errors.As(err, &restErr) {
		restErr = NewInternalServerError(err)
	}

	return &ResponseBody{
		StatusCode: restErr.StatusCode,
		Success:    false,
		Err:        restErr.Err,
		Message:    restErr.Message,
	}
}
