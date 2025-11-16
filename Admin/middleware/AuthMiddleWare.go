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
		
		// Check if the middleware instance or client is nil
		if midwareinstance == nil {
			fmt.Println("AuthMiddleWare instance is nil")
			http.Error(w, "Internal server error: middleware not initialized", http.StatusInternalServerError)
			return
		}
		
		if midwareinstance.Client == nil {
			fmt.Println("Firebase Auth Client is nil")
			http.Error(w, "Internal server error: auth client not initialized", http.StatusInternalServerError)
			return
		}
		
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
            http.Error(w, "Authorization token is required", http.StatusUnauthorized)
            return
        }
// Here we had an issue with this, the problem is that we had to strip out the bearer token. 
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		tokenString = parts[1]
		fmt.Println("Token:", tokenString, "Client:", midwareinstance.Client)
		token, err := midwareinstance.Client.VerifyIDToken(context.Background(), tokenString)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
		fmt.Println("Token:", &token)
		// ctx := context.WithValue(r.Context(), "userClaims", token)
		next.ServeHTTP(w, r)
	})
}