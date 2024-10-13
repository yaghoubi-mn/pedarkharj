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
// @Param number body string true "phone number" example(+98123456789)
// @Param otp body int true "OTP code" example(12345)
// @Param token body string true "Token"
// @Success 200 "Ok. code: code_sent_to_number"
// @Success 303 "Ok. code: go_signup. verify number done. user must signup"
// @Failure 500
// @Failure 400 "BadRequest:<br>code=zero_code_first: Must zero the otp code first.<br>code=wrong_otp: The OTP is wrong.<br>code=number_delay: Wait some minutes.<br>code=invalid_field: a field is invalid"
// @Router /users/verify-number [post]
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

// Signup godoc
// @Summery signup
// @Description Signup user. User must be verify number first.
// @Tags users
// @Accept json
// @Produce json
// @Param number body string true "phone number" example(+98123456789)
// @Param name body string true "full name" example(test)
// @Param token body string true "Token"
// @Param password body string true "Password"
// @Success 200
// @Failure 500
// @Failure 400 "BadRequest:<br>code=verify_number_first: User Must be verify number first<br>code=invalid_field: a field is invalid"<br>code=invalid_token: token is invalid
// @Router /users/signup [post]
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

// Login godoc
// @Summery login user
// @Description login user with number and password
// @Tags users
// @Accept json
// @Produce json
// @Param number body string true "phone number" example(+98123456789)
// @Param password body string true "Password"
// @Success 200
// @Failure 500
// @Failure 400 "BadRequest:<br>code=invalid_field: a field is invalid"
// @Router /users/login [post]
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

// GetUserInfo godoc
// @Description get user info (Authentication Required)
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
// @Failure 500
// @Router /users/info [get]
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

// CheckNumber godoc
// @Description Check number is exist
// @Tags users
// @Accept json
// @Produce json
// @Param number body string true "phone number" example(+98123456789)
// @Success 200
// @Failure 500
// @Failure 400 "BadRequest:<br>code=invalid_field"
// @Router /users/check-number [post]
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

// GetAccessFromRefersh godoc
// @Description Get access token with refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param refresh body string true "refresh"
// @Success 200
// @Failure 500
// @Failure 400 "BadRequest:<br>code=invalid_field"
// @Router /users/refresh [post]
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

// GetAccessFromRefersh godoc
// @Description Choose user avatar
// @Tags users
// @Accept json
// @Produce json
// @Param avatar body string true "Avatar URL"
// @Success 200
// @Failure 500
// @Failure 400 "BadRequest:<br>code=invalid_field"
// @Router /users/avatar [post]
func (h *Handler) ChooseUserAvatar(w http.ResponseWriter, r *http.Request) {

	var avatarInput app_user.AvatarChooseInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&avatarInput)
	defer r.Body.Close()

	iUser := r.Context().Value("user")
	if iUser == nil {
		h.response.ServerErrorResponse(w, errors.New("cannot get user from context"))
		return
	}
	user, ok := iUser.(domain_user.User)
	if !ok {
		h.response.ServerErrorResponse(w, errors.New("cannot cast context user to user"))
		return
	}

	responseDTO := h.appService.ChooseUserAvatar(avatarInput.Avatar, user.ID)
	if responseDTO.UserErr != nil || responseDTO.ServerErr != nil {
		h.response.DTOResponse(w, responseDTO)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}

// GetAccessFromRefersh godoc
// @Description Get list of avatars
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /users/avatar [get]
func (h *Handler) GetAvatars(w http.ResponseWriter, r *http.Request) {

	responseDTO := h.appService.GetAvatars()
	if responseDTO.ServerErr != nil || responseDTO.UserErr != nil {
		h.response.DTOResponse(w, responseDTO)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}
