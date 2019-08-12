package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/l3njo/dropnote-api/controllers"
	"github.com/l3njo/dropnote-api/models"
	u "github.com/l3njo/dropnote-api/utils"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// JwtAuthentication checks validity of the JWT
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{createUser, authUser, getNote, getNotes}
		mayAuth := []string{createNote}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			// TODO
			// notAuth
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}

			// getNote
			prefix := strings.TrimSuffix(value, "{id}")
			suffix := strings.TrimPrefix(requestPath, prefix)
			isUUID := false
			if _, err := uuid.FromString(suffix); err == nil {
				isUUID = true
			}
			if strings.HasPrefix(requestPath, prefix) && isUUID {
				next.ServeHTTP(w, r)
				return
			}

			// codeController
			if strings.Contains(requestPath, "/actions/") || strings.Contains(requestPath, "/assets/") || strings.Contains(requestPath, "/forms/") {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			for _, value := range mayAuth {
				if value == requestPath {
					next.ServeHTTP(w, r)
					return
				}
			}

			resp := u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, resp)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			resp := u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, resp)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			resp := u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, resp)
			return
		}

		if !token.Valid {
			resp := u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, resp)
			return
		}

		ctx := context.WithValue(r.Context(), controllers.UserKey, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
