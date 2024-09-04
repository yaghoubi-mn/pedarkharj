package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

func JSONResponse(w http.ResponseWriter, status int, code int, mapData map[string]interface{}) {
	mapData["code"] = code
	mapData["status"] = status

	json.NewEncoder(w).Encode(mapData)

}

func JSONStructResponse(w http.ResponseWriter, status int, code int, data datatypes.Table) {
	outData := make(map[string]interface{})
	outData["data"] = data
}

// errs example: "name: invalid name"
func JSONErrorResponse(w http.ResponseWriter, status int, code int, errs ...string) {
	if len(errs) == 0 {
		log.Fatalln("errs is required in JSONResponse")
	}

	outData := make(map[string]interface{})
	outData["errors"] = map[string]interface{}{}

	temp := make(map[string]interface{})
	for _, err := range errs {
		splited := strings.Split(err, ":")
		if len(splited) != 2 {
			log.Fatalln("invalid err in JSONErrorResponse")
		}
		temp[splited[0]] = splited[1]
	}

	outData["errors"] = interface{}(temp)

	JSONResponse(w, status, code, outData)

}
