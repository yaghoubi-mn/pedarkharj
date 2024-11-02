package datatypes

import (
	"net/http"
	"time"
)

type Table interface {
	Table() string
}

type CacheRepository interface {
	Save(key string, value map[string]string, expireTime time.Duration) error
	Get(key string) (map[string]string, time.Time, error)
	Delete(key string) error
}

type Validator interface {
	ValidateField(v any, tag string) error
	ValidateFieldByFieldName(fieldName string, fieldValue any, model any) error
}

type Response interface {
	Response(w http.ResponseWriter, status int, code string, mapData Map)
	StructResponse(w http.ResponseWriter, status int, code string, data any)
	ErrorResponse(w http.ResponseWriter, status int, code string, data Map, errs ...error)
	ServerErrorResponse(w http.ResponseWriter, err error)
	DTOErrorResponse(w http.ResponseWriter, responseDTO ResponseDTO)
	InvalidJSONErrorResponse(w http.ResponseWriter, err error)
}
