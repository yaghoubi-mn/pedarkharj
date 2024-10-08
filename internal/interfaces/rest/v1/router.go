package v1

import (
	"net/http"

	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	device_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/device"
	"github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/middleware"
	user_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/user"
)

func NewRouter(userAppService app_user.UserAppService, deviceAppService app_device.DeviceAppService) *http.ServeMux {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()

	// setup json response
	jsonResponse := NewJSONResponse()

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jsonResponse)

	// handlers
	userHandler := user_handler.NewHandler(userAppService, jsonResponse)
	deviceHandler := device_handler.NewHandler(deviceAppService, jsonResponse)

	// user routes
	mux.HandleFunc("POST /users/verify-number", userHandler.VerifyNumber)
	mux.HandleFunc("POST /users/signup", userHandler.SignupUser)
	mux.HandleFunc("POST /users/check-number", userHandler.CheckNumber)
	mux.HandleFunc("POST /users/login", userHandler.Login)
	mux.HandleFunc("POST /users/refresh", userHandler.GetAccessFromRefresh)
	authMux.HandleFunc("GET /users/info", userHandler.GetUserInfo)

	// device routes
	authMux.HandleFunc("POST /devices/logout", deviceHandler.Logout)
	authMux.HandleFunc("POST /devices/logout-all", deviceHandler.LogoutAllUserDevices)

	// handle options
	mux.HandleFunc("OPTIONS /users/verify-number", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("OPTIONS /users/check-number", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("OPTIONS /users/login", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("OPTIONS /users/refresh", func(w http.ResponseWriter, r *http.Request) {})

	// connect muxes
	mux.Handle("/", authMiddleware.EnsureAuthentication(authMux))

	// setup json middleware
	jsonMiddleware := middleware.NewJsonMiddleware(jsonResponse)
	// setup main mux
	m := http.NewServeMux()
	m.Handle("/api/v1/", http.StripPrefix("/api/v1", jsonMiddleware.AddCORSHeaders(jsonMiddleware.EnsureApplicationJson(mux))))

	return m
}
