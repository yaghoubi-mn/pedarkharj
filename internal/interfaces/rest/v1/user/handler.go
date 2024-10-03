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

// VerifyNumber godoc
// @Summery verify number
// @Description verify number with sms
// @Tags users
// @Accept json
// @Produce json
// @Parm number code token
// @Success 200
func (h *Handler) VerifyNumber(w http.ResponseWriter, r *http.Request) {
	var verifyNumberInput app_user.VerifyNumberInput
	// decode body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&verifyNumberInput)

	defer r.Body.Close()

	userAgent := utils.GetUserAgent(r)
	userIP := utils.GetIPAddress(r)

	step, responseDTO := h.appService.VerifyNumber(verifyNumberInput, userAgent, userIP)
	if responseDTO.ServerErr != nil {
		h.response.ServerErrorResponse(w, responseDTO.ServerErr)
		return
	}

	if responseDTO.UserErr != nil {
		h.response.ErrorResponse(w, http.StatusBadRequest, responseDTO.ResponseCode, responseDTO.UserErr)
		return
	}

	// otp code sent to number
	if step == 1 {
		responseDTO.Data["msg"] = "Code sent to number"
		h.response.Response(w, http.StatusOK, responseDTO.ResponseCode, responseDTO.Data)
		return
	}

	// user sent otp code and otp is currect
	if step == 2 {
		responseDTO.Data["msg"] = "Number verified. Go signup"
		h.response.Response(w, 303, responseDTO.ResponseCode, responseDTO.Data)
		return
	}

	// user sent otp code and otp is currect. user already exists in database
	if step == 3 {
		responseDTO.Data["msg"] = "You are in!"
		h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
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

	responseDTO := h.appService.Signup(userInput, userAgent, userIP)
	if responseDTO.ServerErr != nil {
		h.response.ServerErrorResponse(w, responseDTO.ServerErr)
		return
	}

	if responseDTO.UserErr != nil {
		h.response.ErrorResponse(w, http.StatusBadRequest, responseDTO.ResponseCode, responseDTO.UserErr)
		return

	}

	h.response.Response(w, 200, "", responseDTO.Data)
}

// login user with number and password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userInput app_user.LoginUserInput
	// decode body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&userInput)

	responseDTO := h.appService.Login(userInput, utils.GetUserAgent(r), utils.GetIPAddress(r))
	if responseDTO.ServerErr != nil {
		h.response.ServerErrorResponse(w, responseDTO.ServerErr)
		return
	}
	if responseDTO.UserErr != nil {
		h.response.ErrorResponse(w, 400, responseDTO.ResponseCode, responseDTO.UserErr)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
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

	responseDTO := h.appService.GetUserInfo(user)

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}

func (h *Handler) CheckNumber(w http.ResponseWriter, r *http.Request) {

	var numberInput app_user.NumberInput

	json.NewDecoder(r.Body).Decode(&numberInput)

	responseDTO := h.appService.CheckNumber(numberInput)

	if responseDTO.ServerErr != nil {
		h.response.ServerErrorResponse(w, responseDTO.ServerErr)
		return
	}
	if responseDTO.UserErr != nil {
		h.response.ErrorResponse(w, 400, responseDTO.ResponseCode, responseDTO.UserErr)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}

func (h *Handler) GetAccessFromRefresh(w http.ResponseWriter, r *http.Request) {

	var refreshInput app_user.RefreshInput

	json.NewDecoder(r.Body).Decode(&refreshInput)

	responseDTO := h.appService.GetAccessFromRefresh(refreshInput.Refresh)
	if responseDTO.ServerErr != nil {
		h.response.ServerErrorResponse(w, responseDTO.ServerErr)
		return
	}
	if responseDTO.UserErr != nil {
		h.response.ErrorResponse(w, 400, responseDTO.ResponseCode, responseDTO.UserErr)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}
