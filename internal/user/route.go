package user

import "net/http"

func Route(mux *http.ServeMux, userHandler Handler) {
	userMux := http.NewServeMux()
	userMux.HandleFunc("signup", userHandler.CreateUser)

	// mux.Handle("/user/", middleware.authMiddleware(userMux))
	mux.Handle("/user/", userMux)
}
