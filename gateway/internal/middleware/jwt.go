package middleware

import (
	"fmt"
	config "github.com/Nixonxp/discord/gateway/configs"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

var secretKey *string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip if auth routes
		if strings.Contains(strings.ToLower(r.URL.Path), "auth") {
			next.ServeHTTP(w, r)
		} else {
			if secretKey == nil {
				cfg := config.GetConfig()
				secretKey = &cfg.Application.AuthSecretKey
				fmt.Println("get key from config")
			}
			authorizationHeader := r.Header.Get("Authorization")
			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
			token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(*secretKey), nil
			})
			if err != nil {
				http.Error(w, "Invalid or expired JWT token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(*jwt.StandardClaims)
			userID := claims.Id

			ctx := metadata.AppendToOutgoingContext(r.Context(), "userId", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
