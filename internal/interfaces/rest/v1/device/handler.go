package device_handler

import (
	"errors"
	"net/http"

	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/interfaces_rest_v1_shared"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type Handler struct {
	appService app_device.DeviceAppService
	response   interfaces_rest_v1_shared.Response
}

func NewHandler(appService app_device.DeviceAppService, response interfaces_rest_v1_shared.Response) Handler {
	return Handler{
		appService: appService,
		response:   response,
	}
}

// Logout godoc
// @Description logout current user device
// @Tags devices
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /devices/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	iUser := r.Context().Value("user")
	if iUser == nil {
		h.response.ServerErrorResponse(w, errors.New("cannot get user from context"))
		return
	}

	user, ok := iUser.(domain_user.User)
	if !ok {
		h.response.ServerErrorResponse(w, errors.New("cannot cast to user"))
		return
	}

	responseDTO := h.appService.Logout(user.ID, utils.GetUserAgent(r))

	if responseDTO.ServerErr != nil || responseDTO.UserErr != nil {
		h.response.DTOErrorResponse(w, responseDTO)
		return
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}

// LogoutAllUserDevices godoc
// @Description logout all user devices
// @Tags devices
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /devices/logout-all [post]
func (h *Handler) LogoutAllUserDevices(w http.ResponseWriter, r *http.Request) {

	iUser := r.Context().Value("user")
	if iUser == nil {
		h.response.ServerErrorResponse(w, errors.New("cannot get user from context"))
		return
	}

	user, ok := iUser.(domain_user.User)
	if !ok {
		h.response.ServerErrorResponse(w, errors.New("cannot cast to user"))
		return
	}

	responseDTO := h.appService.LogoutAllUserDevices(user.ID)
	if responseDTO.ServerErr != nil || responseDTO.UserErr != nil {
		h.response.DTOErrorResponse(w, responseDTO)
	}

	h.response.Response(w, 200, responseDTO.ResponseCode, responseDTO.Data)
}
