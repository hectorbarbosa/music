// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// InfoGet implements GET /info operation.
//
// GET /info
func (UnimplementedHandler) InfoGet(ctx context.Context, params InfoGetParams) (r InfoGetRes, _ error) {
	return r, ht.ErrNotImplemented
}
