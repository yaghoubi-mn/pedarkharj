package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

func JSONResponse(w http.ResponseWriter, status int, code string, mapData datatypes.Map) {
	mapData["code"] = code
	mapData["status"] = status

	json.NewEncoder(w).Encode(mapData)

}

func JSONStructResponse(w http.ResponseWriter, status int, code string, data datatypes.Table) {
	outData := make(datatypes.Map)
	outData["data"] = data
}

// errs example: "name: invalid name"
func JSONErrorResponse(w http.ResponseWriter, status int, code string, errMap map[string]string) {
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
	log.Println(err.Error())
	JSONResponse(w, http.StatusInternalServerError, "", datatypes.Map{"msg": "Server error"})
}
