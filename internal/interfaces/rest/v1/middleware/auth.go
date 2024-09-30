package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type authMiddleware struct {
	response datatypes.Response
}

func NewAuthMiddleware(response datatypes.Response) authMiddleware {
	return authMiddleware{
		response: response,
	}
}

func (a *authMiddleware) EnsureAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access := r.Header.Get("Authorization")
		if access == "" {
			a.response.ErrorResponse(w, 401, rcodes.Unauthenticated, errors.New("authentication is required"))
			return
		}

		if strings.Index(access, "Bearer ") != 0 || len(access) < 8 {
			a.response.ErrorResponse(w, 400, rcodes.InvalidHeader, errors.New("invalid Authorization header format"))
			return
		}

		access = access[7:]
		fmt.Println(access)

		var user domain_user.User
		var err error
		user.ID, user.Name, user.Number, user.IsRegistered, err = jwt.GetUserFromAccess(access)
		if err != nil {
			fmt.Println(err)
			a.response.ErrorResponse(w, 401, rcodes.InvalidToken, errors.New("invalid token"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next.ServeHTTP(w, r)
	})
}
