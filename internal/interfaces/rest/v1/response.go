package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

type jsonResponse struct {
}

func NewJSONResponse() datatypes.Response {
	return &jsonResponse{}
}

func (j *jsonResponse) Response(w http.ResponseWriter, status int, code string, mapData datatypes.Map) {
	mapData["code"] = code
	mapData["status"] = status

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(mapData)

	slog.Info("request info",
		slog.Int("status", status),
		slog.Any("code", code),
		slog.Any("data", mapData),
	)
}

func (j *jsonResponse) StructResponse(w http.ResponseWriter, status int, code string, data any) {
	outData := make(datatypes.Map)
	outData["data"] = data
	j.Response(w, status, code, outData)
}

// errs example: "name: invalid name"
func (j *jsonResponse) ErrorResponse(w http.ResponseWriter, status int, code string, errs ...error) {
	if errs == nil {
		slog.Error("err is required in JSONErrorResponse")
	}

	outData := make(datatypes.Map)
	outData["errors"] = datatypes.Map{}

	temp := make(map[string]string)

	for _, err := range errs {
		splited := strings.Split(err.Error(), ": ")
		if len(splited) == 1 {
			temp["non-field"] = splited[0]
		} else {
			temp[splited[0]] = splited[1]
		}
	}

	outData["errors"] = interface{}(temp)

	j.Response(w, status, code, outData)

}

func (j *jsonResponse) ServerErrorResponse(w http.ResponseWriter, err error) {
	slog.Error("SERVER ERROR", "error", err.Error())
	j.Response(w, http.StatusInternalServerError, "", datatypes.Map{"msg": "Server error"})
}
