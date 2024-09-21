package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

func JSONResponse(w http.ResponseWriter, status int, code rcodes.ResponseCode, mapData datatypes.Map) {
	mapData["code"] = code
	mapData["status"] = status

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(mapData)

	log.Println(status, code, mapData)
}

func JSONStructResponse(w http.ResponseWriter, status int, code rcodes.ResponseCode, data datatypes.Table) {
	outData := make(datatypes.Map)
	outData["data"] = data
}

// errs example: "name: invalid name"
func JSONErrorResponse(w http.ResponseWriter, status int, code rcodes.ResponseCode, errMap map[string]string) {
	if errMap == nil {
		log.Fatalln("errMap is required in JSONResponse")
	}

	outData := make(datatypes.Map)
	outData["errors"] = datatypes.Map{}

	temp := make(datatypes.Map)
	for key := range errMap {
		temp[key] = errMap[key]
	}

	outData["errors"] = interface{}(temp)

	JSONResponse(w, status, code, outData)

}

func JSONServerError(w http.ResponseWriter, err error) {
	log.Println("SERVER ERROR: ", err.Error())
	JSONResponse(w, http.StatusInternalServerError, "", datatypes.Map{"msg": "Server error"})
}
