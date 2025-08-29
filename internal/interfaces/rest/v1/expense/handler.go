package expense_handler

import (
	"encoding/json"
	"errors"
	"net/http"

	app_expense "github.com/yaghoubi-mn/pedarkharj/internal/application/expense"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	interfaces_rest_v1_shared "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/shared"
)

type Handler struct {
	appService app_expense.ExpenseAppService
	response   interfaces_rest_v1_shared.Response
}

func NewHandler(appService app_expense.ExpenseAppService, response interfaces_rest_v1_shared.Response) Handler {
	return Handler{
		appService: appService,
		response:   response,
	}
}

// SendOTP godoc
// @Summery create new expense
// @Description create new expense.
// @Tags expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name body string true "expense name"
// @Param description body string true "expense description"
// @Param creditors body map[string]uint64 true "creditors key value list: phone number is key and credit amount is value" exmaple("{"+989123456789": 2000, "+989123456788": 5000}")
// @Param debtors body []string true "list of debtors phone number" example("["+989123456786", "+989123456787"]")
// @Success 200 "Ok"
// @Failure 500
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 "BadRequest:<br>code=invalid_field: a field is invalid"
// @Router /expenses [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var input app_expense.ExpenseInputWithPhoneNumber
	// decode body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&input)
	defer r.Body.Close()

	if err != nil {
		h.response.InvalidJSONErrorResponse(w, err)
		return
	}

	iUser := r.Context().Value("user")
	if iUser == nil {
		h.response.ServerErrorResponse(w, errors.New("user is nil in request context"))
		return
	}

	user, ok := iUser.(app_user.JWTUser)
	if !ok {
		h.response.ServerErrorResponse(w, errors.New("cannot cast request context user"))
		return
	}

	responseDTO := h.appService.Create(input, user.ID, user.PhoneNumber)
	if responseDTO.ServerErr != nil || responseDTO.UserErr != nil {
		h.response.DTOErrorResponse(w, responseDTO)
		return
	}

	h.response.Response(w, http.StatusOK, responseDTO.ResponseCode, responseDTO.Data)

}
