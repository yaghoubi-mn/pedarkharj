package v1

import (
	"net/http"

	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/middleware"
	user_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/user"
)

func NewRouter(userAppService app_user.UserAppService) *http.ServeMux {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()

	// setup json response
	jsonResponse := NewJSONResponse()

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jsonResponse)

	// handlers
	userHandler := user_handler.NewHandler(userAppService, jsonResponse)

	// user routes
	mux.HandleFunc("POST /users/verify-number", userHandler.VerifyNumber)
	mux.HandleFunc("POST /users/signup", userHandler.SignupUser)
	authMux.HandleFunc("GET /users/info", userHandler.GetUserInfo)

	// connect muxes
	mux.Handle("/", authMiddleware.EnsureAuthentication(authMux))

	// setup json middleware
	jsonMiddleware := middleware.NewJsonMiddleware(jsonResponse)
	// setup main mux
	m := http.NewServeMux()
	m.Handle("/api/v1/", http.StripPrefix("/api/v1", jsonMiddleware.EnsureApplicationJson(mux)))

	return m
}
