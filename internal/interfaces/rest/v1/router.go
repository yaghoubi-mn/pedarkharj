package v1

import (
	"encoding/json"
	"net/http"

	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	device_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/device"
	"github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/middleware"
	user_handler "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/user"
)

func NewRouter(userAppService app_user.UserAppService, deviceAppService app_device.DeviceAppService) *http.ServeMux {
	mux := http.NewServeMux()
	// authMux := http.NewServeMux()

	// setup json response
	jsonResponse := NewJSONResponse()

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jsonResponse)

	// handlers
	userHandler := user_handler.NewHandler(userAppService, jsonResponse)
	deviceHandler := device_handler.NewHandler(deviceAppService, jsonResponse)

	// handle 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			data := make(map[string]any)
			data["status"] = 404
			data["msg"] = "page not found"
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		}
	})

	// user routes
	registerRouteFunc(mux, "POST", "/users/verify-number", userHandler.VerifyNumber)
	registerRouteFunc(mux, "POST", "/users/signup", userHandler.SignupUser)
	registerRouteFunc(mux, "POST", "/users/check-number", userHandler.CheckNumber)
	registerRouteFunc(mux, "POST", "/users/login", userHandler.Login)
	registerRouteFunc(mux, "POST", "/users/refresh", userHandler.GetAccessFromRefresh)
	registerRoute(mux, "GET", "/users/info", authMiddleware.EnsureAuthentication(http.HandlerFunc(userHandler.GetUserInfo)))

	// device routes
	registerRoute(mux, "POST", "/devices/logout", authMiddleware.EnsureAuthentication(http.HandlerFunc((deviceHandler.Logout))))
	registerRoute(mux, "POST", "/devices/logout-all", authMiddleware.EnsureAuthentication(http.HandlerFunc(deviceHandler.LogoutAllUserDevices)))

	// connect muxes
	// mux.Handle("/", authMiddleware.EnsureAuthentication(authMux))

	// setup json middleware
	jsonMiddleware := middleware.NewJsonMiddleware(jsonResponse)
	// setup main mux
	m := http.NewServeMux()
	m.Handle("/api/v1/", http.StripPrefix("/api/v1", jsonMiddleware.AddCORSHeaders(jsonMiddleware.EnsureApplicationJson(mux))))

	return m
}

func registerRouteFunc(mux *http.ServeMux, method string, url string, handler func(http.ResponseWriter, *http.Request)) {

	registerRoute(mux, method, url, http.HandlerFunc(handler))

}

// func registerRouteFuncHelper(next http.Handler) http.Handler {
// 	return
// }

func registerRoute(mux *http.ServeMux, method string, url string, handle http.Handler) {
	mux.Handle(method+" "+url, handle)

	mux.HandleFunc("OPTIONS "+url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "content-type, access-control-allow-origin, accept, user-agent, authorization")
		w.Header().Set("Access-Control-Allow-Max-Age", "86400")
	})
}
