package user

import (
	"encoding/json"
	"errors"
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

func (h *Handler) VerifyNumber(w http.ResponseWriter, r *http.Request) {
	var verifyNumberInput VerifyNumberInput
	json.NewDecoder(r.Body).Decode(&verifyNumberInput)
	defer r.Body.Close()

	step, code, tokens, errMap, err := h.service.VerifyNumber(verifyNumberInput)
	if err != nil {
		utils.JSONServerError(w, err)
		return
	}

	if errMap != nil {
		utils.JSONErrorResponse(w, http.StatusBadRequest, code, errMap)
		return
	}

	// otp code sent to number
	if step == 1 {
		utils.JSONResponse(w, http.StatusOK, code, datatypes.Map{"msg": "Code sent to number.", "token": tokens["token"]})
		return
	}

	// user sent otp code and otp is currect
	if step == 2 {
		utils.JSONResponse(w, 303, code, datatypes.Map{"msg": "Number verified. Go signup."})
		return
	}

	// user sent otp code and otp is currect. user already exists in database
	if step == 3 {
		utils.JSONResponse(w, 200, code, datatypes.Map{"msg": "You are in!", "refresh": tokens["refresh"], "access": tokens["access"]})
		return
	}

	utils.JSONServerError(w, errors.New("unhandled state"))

}

func (h *Handler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var userInput SignupUserInput
	json.NewDecoder(r.Body).Decode(&userInput)
	defer r.Body.Close()

	tokens, code, errMap, err := h.service.Signup(userInput)
	if err != nil {
		utils.JSONServerError(w, err)
		return
	}

	if errMap != nil {
		utils.JSONErrorResponse(w, http.StatusBadRequest, code, errMap)
		return

	}

	utils.JSONResponse(w, 200, "", datatypes.Map{"msg": "done", "refresh": tokens["refresh"], "access": tokens["access"]})
}
