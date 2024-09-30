package datatypes

import (
	"net/http"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type Table interface {
	Table() string
}

type CacheRepository interface {
	Save(key string, value string, expireTime time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Validator interface {
	ValidateField(v any, tag string) error
}

type Response interface {
	Response(w http.ResponseWriter, status int, code rcodes.ResponseCode, mapData Map)
	StructResponse(w http.ResponseWriter, status int, code rcodes.ResponseCode, data any)
	ErrorResponse(w http.ResponseWriter, status int, code rcodes.ResponseCode, errs ...error)
	ServerErrorResponse(w http.ResponseWriter, err error)
}
