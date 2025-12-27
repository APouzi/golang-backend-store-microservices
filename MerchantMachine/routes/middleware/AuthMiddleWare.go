package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type AuthMiddleWare struct {
	Client *auth.Client
}


func(midwareinstance *AuthMiddleWare) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		tokenString = parts[1]

		fmt.Println(tokenString, "and the client!", midwareinstance.Client)
		token, err := midwareinstance.Client.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		fmt.Println("Token:", &token)

		// pull out email from token claims and print it out
		var email string
		if e, ok := token.Claims["email"].(string); ok && e != "" {
			email = e
		} else if e2, ok := token.Claims["user_email"].(string); ok && e2 != "" {
			// some tokens might use a different claim key
			email = e2
		}
		fmt.Println("Email:", email)

		// attach email to request context for downstream handlers
		ctx := context.WithValue(r.Context(), "userEmail", email)
		r = r.WithContext(ctx)

		// ctx := context.WithValue(r.Context(), "userClaims", token)
		next.ServeHTTP(w, r)
	})
}

func(midwareinstance *AuthMiddleWare) HelloAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from AuthMiddleware!")
		next.ServeHTTP(w, r)
	})
}