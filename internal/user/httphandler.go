package user

import (
	"encoding/json"
	"net/http"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type Handler struct {
	service UserService
}

func NewHandler(service UserService) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	err := h.service.CreateUser(user)
	if err != nil {
		utils.JSONErrorResponse(w, 400, 10001, err.Error())
		return
	}

	utils.JSONResponse(w, 200, 0, datatypes.Map{"msg": "done"})
	return
}
