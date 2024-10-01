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

		if r.Method == "OPTIIONS" {

			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "content-type, access-control-allow-origin, accept, user-agent, authorization")
			w.Header().Set("Access-Control-Allow-Max-Age", "86400")
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			j.response.ErrorResponse(w, 400, rcodes.InvalidHeader, errors.New("header application/json is required"))
			return
		}
		// w.Header().Add("")
		next.ServeHTTP(w, r)

	})
}
