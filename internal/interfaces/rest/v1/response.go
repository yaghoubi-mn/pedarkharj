package v1

import (
	"encoding/json"
	"errors"
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
func (j *jsonResponse) ErrorResponse(w http.ResponseWriter, status int, code string, data datatypes.Map, errs ...error) {
	if errs == nil {
		slog.Error("err is required in JSONErrorResponse")
	}

	if data == nil {
		data = make(datatypes.Map)
	}

	data["errors"] = datatypes.Map{}

	temp := make(map[string]string)

	for _, err := range errs {
		splited := strings.Split(err.Error(), ": ")
		if len(splited) == 1 {
			temp["non_field"] = splited[0]
		} else {
			temp[splited[0]] = splited[1]
		}
	}

	data["errors"] = interface{}(temp)

	j.Response(w, status, code, data)

}

func (j *jsonResponse) ServerErrorResponse(w http.ResponseWriter, err error) {
	slog.Error("SERVER ERROR", "error", err.Error())
	j.Response(w, http.StatusInternalServerError, "", datatypes.Map{"msg": "Server error"})
}

// check ServerErr and UserErr
func (j *jsonResponse) DTOErrorResponse(w http.ResponseWriter, responseDTO datatypes.ResponseDTO) {

	if responseDTO.ServerErr != nil {
		j.ServerErrorResponse(w, responseDTO.ServerErr)
	} else if responseDTO.UserErr != nil {
		j.ErrorResponse(w, 400, responseDTO.ResponseCode, responseDTO.Data, responseDTO.UserErr)
	}

}

func (j *jsonResponse) InvalidJSONErrorResponse(w http.ResponseWriter, err error) {

	j.ErrorResponse(w, 400, "invalid_json", nil, errors.New("invalid json"))
}
