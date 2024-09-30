package user

import (
	"encoding/json"
	"errors"
	"net/http"

	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type Handler struct {
	appService app_user.UserAppService
	response   datatypes.Response
}

func NewHandler(appService app_user.UserAppService, response datatypes.Response) Handler {
	return Handler{
		appService: appService,
		response:   response,
	}
}

func (h *Handler) VerifyNumber(w http.ResponseWriter, r *http.Request) {
	var verifyNumberInput app_user.VerifyNumberInput
	// decode body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&verifyNumberInput)

	defer r.Body.Close()

	userAgent := utils.GetUserAgent(r)
	userIP := utils.GetIPAddress(r)

	step, code, tokens, userErr, serverErr := h.appService.VerifyNumber(verifyNumberInput, userAgent, userIP)
	if serverErr != nil {
		h.response.ServerErrorResponse(w, serverErr)
		return
	}

	if userErr != nil {
		h.response.ErrorResponse(w, http.StatusBadRequest, code, userErr)
		return
	}

	// otp code sent to number
	if step == 1 {
		h.response.Response(w, http.StatusOK, code, datatypes.Map{"msg": "Code sent to number.", "token": tokens["token"]})
		return
	}

	// user sent otp code and otp is currect
	if step == 2 {
		h.response.Response(w, 303, code, datatypes.Map{"msg": "Number verified. Go signup."})
		return
	}

	// user sent otp code and otp is currect. user already exists in database
	if step == 3 {
		h.response.Response(w, 200, code, datatypes.Map{"msg": "You are in!", "refresh": tokens["refresh"], "access": tokens["access"], "accessExpireSeconds": tokens["accessExpireSeconds"]})
		return
	}

	h.response.ServerErrorResponse(w, errors.New("unhandled state"))

}

func (h *Handler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var userInput app_user.SignupUserInput
	// decode body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&userInput)

	defer r.Body.Close()

	userAgent := utils.GetUserAgent(r)
	userIP := utils.GetIPAddress(r)

	tokens, code, userErr, serverErr := h.appService.Signup(userInput, userAgent, userIP)
	if serverErr != nil {
		h.response.ServerErrorResponse(w, serverErr)
		return
	}

	if userErr != nil {
		h.response.ErrorResponse(w, http.StatusBadRequest, code, userErr)
		return

	}

	h.response.Response(w, 200, "", datatypes.Map{"msg": "done", "refresh": tokens["refresh"], "access": tokens["access"], "accessExpireSeconds": tokens["accessExpireSeconds"]})
}

// login user with number and password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {

	iUser := r.Context().Value("user")
	if iUser == nil {
		h.response.ServerErrorResponse(w, errors.New("user is nil in request context"))
		return
	}

	user, ok := iUser.(domain_user.User)
	if !ok {
		h.response.ServerErrorResponse(w, errors.New("cannot cast request context user"))
		return
	}

	outData := h.appService.GetUserInfo(user)

	h.response.StructResponse(w, 200, "", outData)
}
