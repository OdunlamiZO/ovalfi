package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/odunlamizo/ovalfi/internal/model"
	"github.com/odunlamizo/ovalfi/internal/util"
)

// authorize user by validating jwt
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorMessage := model.ResponseMessage{Message: "Authorization token not provided"}
			responseBody, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBody)
			return
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			errorMessage := model.ResponseMessage{Message: "Invalid authorization header"}
			responseBody, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBody)
			return
		}

		claims, err := util.ValidateToken(authHeaderParts[1])
		if err != nil {
			errorMessage := model.ResponseMessage{Message: "Error validating token"}
			responseBody, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBody)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims["sub"]) // to keep track of authenticated user
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
