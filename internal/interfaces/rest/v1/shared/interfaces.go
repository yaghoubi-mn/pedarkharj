package interfaces_rest_v1_shared

import (
	"net/http"

	app_shared "github.com/yaghoubi-mn/pedarkharj/internal/application/shared"
)

type Response interface {
	Response(w http.ResponseWriter, status int, code string, mapData map[string]any)
	StructResponse(w http.ResponseWriter, status int, code string, data any)
	ErrorResponse(w http.ResponseWriter, status int, code string, data map[string]any, errs ...error)
	ServerErrorResponse(w http.ResponseWriter, err error)
	DTOErrorResponse(w http.ResponseWriter, responseDTO app_shared.ResponseDTO)
	InvalidJSONErrorResponse(w http.ResponseWriter, err error)
}
