package v1

import (
	"net/http"

	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	user_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/user"
)

func NewRouter(userAppService app_user.UserAppService) *http.ServeMux {
	mux := http.NewServeMux()

	jsonResponse := NewJSONResponse()

	// handlers
	userHandler := user_handler.NewHandler(userAppService, jsonResponse)

	// user routes
	mux.HandleFunc("POST /users/verify-number", userHandler.VerifyNumber)
	mux.HandleFunc("POST /users/signup", userHandler.SignupUser)

	// mux.Handle("/", middleware(fun))
	m := http.NewServeMux()

	m.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	return m
}
