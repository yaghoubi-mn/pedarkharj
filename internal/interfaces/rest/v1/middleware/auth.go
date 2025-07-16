package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1/shared"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type authMiddleware struct {
	response interfaces_rest_v1_shared.Response
}

func NewAuthMiddleware(response interfaces_rest_v1_shared.Response) authMiddleware {
	return authMiddleware{
		response: response,
	}
}

func (a *authMiddleware) EnsureAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access := r.Header.Get("Authorization")
		if access == "" {
			a.response.ErrorResponse(w, 401, rcodes.Unauthenticated, nil, errors.New("authentication is required"))
			return
		}

		if strings.Index(access, "Bearer ") != 0 || len(access) < 8 {
			a.response.ErrorResponse(w, 400, rcodes.InvalidHeader, nil, errors.New("invalid Authorization header format"))
			return
		}

		access = access[7:]

		var user app_user.JWTUser
		var err error
		user.ID, user.Name, user.Number, user.IsRegistered, err = jwt.GetUserFromAccess(access)
		if err != nil {
			fmt.Println("JWT ERROR:", err)
			a.response.ErrorResponse(w, 401, rcodes.InvalidToken, nil, errors.New("authorization: invalid token"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next.ServeHTTP(w, r)
	})
}
