package user

import "net/http"

func Route(prefix string, mux *http.ServeMux, userHandler Handler) {

	// mux := http.NewServeMux()
	mux.HandleFunc("POST "+prefix+"/signup", userHandler.SignupUser)
	mux.HandleFunc("POST "+prefix+"/verify-number", userHandler.VerifyNumber)
	// mux.Handle("/user/", middleware.authMiddleware(userMux))

	// mux.Handle("/users/*", userMux)
}
