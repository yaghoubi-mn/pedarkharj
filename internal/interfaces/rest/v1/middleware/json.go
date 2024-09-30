package middleware

import (
	"errors"
	"net/http"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type jsonMiddleware struct {
	response datatypes.Response
}

func NewJsonMiddleware(response datatypes.Response) jsonMiddleware {
	return jsonMiddleware{
		response: response,
	}
}

func (j *jsonMiddleware) EnsureApplicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			j.response.ErrorResponse(w, 400, rcodes.InvalidHeader, errors.New("header application/json is required"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
