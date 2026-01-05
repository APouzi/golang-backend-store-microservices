package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/APouzi/MerchantMachinee/routes/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

//You can also authenticate with Firebase using a Google Account by handling the sign-in flow with the Sign In With Google library:
// https://firebase.google.com/docs/auth/web/google-signin#expandable-2

type JWTtest struct{
	Token string `json:"JWT"`
}

type AuthMiddleWareStruct struct{
	// db *sql.DB
	firebaseApp *firebase.App
	validationclient *auth.Client
	redisClient      *redis.Client
}

func InjectSystemRefrences(firebaseApp *firebase.App, redisClient *redis.Client) *AuthMiddleWareStruct{
	client, err := firebaseApp.Auth(context.Background())
	if err != nil {
		fmt.Println("Failed to initialize Firebase Auth client:", err)
		return nil // or handle error appropriately
	}

	authMiddleWareInstance := AuthMiddleWareStruct{
		firebaseApp:      firebaseApp,
		validationclient: client,
		redisClient:      redisClient,
	}
	return &authMiddleWareInstance
}

//Initialize the SDK in non-Google environments: If you are working in a non-Google server environment (This appp, I believe) in which default credentials lookup can't be fully automated, you can initialize the SDK with an exported service account key file. 
// https://firebase.google.com/docs/admin/setup#initialize_the_sdk_in_non-google_environments
func(db *AuthMiddleWareStruct) ValidateToken(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		jwttoken := r.Header.Get("Authorization")
		// When the jwttoken comes in, it will input "bearer" into the token and we have to remove this from the token so we can parse it. 
		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
			return []byte("Testing key"), nil
		})
		if err != nil{
			fmt.Println("ValidateToken Failed")
			helpers.ErrorJSON(w,err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

type ctxKey string

const ctxUserEmail ctxKey = "userEmail"
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ctxUserEmail).(string)
	return email, ok
}

// CheckUserRegistration is an HTTP middleware that validates a Firebase OAuth (ID) token and injects the
// authenticated user's email into the request context.
//
// The returned http.Handler performs the following steps:
// - Verifies the middleware instance and its validation client are initialized; if not, responds with 500.
// - Reads the Authorization header and requires the "Bearer {token}" format; missing or malformed headers
//   result in a 401 Unauthorized response.
// - Emits concise debug output about the token (length and prefix/suffix) without logging the full token.
// - Performs a basic syntactic check that the token looks like a JWT (expects exactly two '.' separators);
//   if the format is invalid, responds 401.
// - Calls midwareinstance.validationclient.VerifyIDToken(ctx, tokenString) to validate the Firebase ID token;
//   any verification error is returned as a 401 Unauthorized response.
// - Extracts the user's email from the verified token's claims (first trying "email", then "user_email").
// - Stores the extracted email in the request context under the key "userEmail" and forwards the request to the
//   next handler. If authentication fails at any point, the middleware writes an appropriate HTTP error and
//   does not call the next handler.
func(midwareinstance *AuthMiddleWareStruct) CheckUserRegistration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow OPTIONS requests to pass through (for CORS preflight)
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("We have entere  check")

		if midwareinstance == nil || midwareinstance.validationclient == nil {
			helpers.ErrorJSON(w, fmt.Errorf("auth client not initialized"), http.StatusInternalServerError)
			return
		}

		authz := strings.TrimSpace(r.Header.Get("Authorization"))

		// If Authorization header is missing, check the body for "token" field
		if authz == "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				// Restore the body for downstream handlers
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				var bodyToken struct {
					Token string `json:"token"`
				}
				if json.Unmarshal(bodyBytes, &bodyToken) == nil && bodyToken.Token != "" {
					authz = "Bearer " + bodyToken.Token
					// Inject into header so downstream handlers (like RegisterCustomer) can find it
					r.Header.Set("Authorization", authz)
				}
			}
		}
		if authz == "" {
			helpers.ErrorJSON(w, fmt.Errorf("authorization token is required"), http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			helpers.ErrorJSON(w, fmt.Errorf("authorization header format must be Bearer {token}"), http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimSpace(authz[len("Bearer "):])
		if strings.Count(tokenString, ".") != 2 {
			helpers.ErrorJSON(w, fmt.Errorf("invalid token format (expected JWT)"), http.StatusUnauthorized)
			return
		}

		token, err := midwareinstance.validationclient.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			fmt.Printf("Error verifying token: %v\n", err)
			helpers.ErrorJSON(w, fmt.Errorf("invalid token: %v", err), http.StatusUnauthorized)
			return
		} 

		email, _ := token.Claims["email"].(string)
		if email == "" {
			email, _ = token.Claims["user_email"].(string)
		}
		if email == "" {
			helpers.ErrorJSON(w, fmt.Errorf("email claim missing"), http.StatusUnauthorized)
			return
		}

		exists, err := midwareinstance.redisClient.Get(r.Context(), email).Result()
		if err != nil && err != redis.Nil {
			helpers.ErrorJSON(w, fmt.Errorf("redis error: %v", err), http.StatusInternalServerError)
			return
		}
		if exists == "" {
			// Use environment variable for service URL, defaulting to the docker service name
			dbURL := os.Getenv("DBLAYER_URL")
			if dbURL == "" {
				dbURL = "http://dblayer:8080"
			}
			
			// Prepare payload
			payload := CreateUserWithProfileRequest{
				Email: email,
			}

			// 1. Try to fill from Token Claims
			if name, ok := token.Claims["name"].(string); ok {
				parts := strings.Fields(name)
				if len(parts) > 0 {
					payload.FirstName = parts[0]
				}
				if len(parts) > 1 {
					payload.LastName = strings.Join(parts[1:], " ")
				}
			}
			if phone, ok := token.Claims["phone_number"].(string); ok {
				payload.PhoneNumberMobileE164 = phone
			}

			// 2. Try to fill from Request Body (if it's a registration request)
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore immediately
				
				var bodyRequest CreateUserWithProfileRequest
				if json.Unmarshal(bodyBytes, &bodyRequest) == nil {
					// Overwrite/Merge fields if present in body
					if bodyRequest.FirstName != "" { payload.FirstName = bodyRequest.FirstName }
					if bodyRequest.LastName != "" { payload.LastName = bodyRequest.LastName }
					if bodyRequest.PhoneNumberMobileE164 != "" { payload.PhoneNumberMobileE164 = bodyRequest.PhoneNumberMobileE164 }
					if bodyRequest.PhoneNumberHomeE164 != "" { payload.PhoneNumberHomeE164 = bodyRequest.PhoneNumberHomeE164 }
					if bodyRequest.PrimaryShippingAddressID != nil { payload.PrimaryShippingAddressID = bodyRequest.PrimaryShippingAddressID }
					if bodyRequest.PrimaryBillingAddressID != nil { payload.PrimaryBillingAddressID = bodyRequest.PrimaryBillingAddressID }
					if bodyRequest.PreferredLocale != nil { payload.PreferredLocale = bodyRequest.PreferredLocale }
					if bodyRequest.PreferredTimeZone != nil { payload.PreferredTimeZone = bodyRequest.PreferredTimeZone }
				}
			}
			jsonData, _ := json.Marshal(payload)

			// Create request with context to propagate timeouts/cancellations
			req, err := http.NewRequestWithContext(r.Context(), "POST", dbURL+"/users/profile", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("Error creating profile sync request: %v\n", err)
			} else {
				req.Header.Set("Content-Type", "application/json")
				
				// Execute request
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error syncing user profile: %v\n", err)
				} else {
					if resp.StatusCode >= 400 {
						bodyBytes, _ := io.ReadAll(resp.Body)
						fmt.Printf("Error syncing user profile. Status: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
					} else {
						fmt.Println("User profile synced successfully")
					}
					// Important: Close the body to prevent resource leaks
					resp.Body.Close()
				}
			}
		}



		ctx := context.WithValue(r.Context(), ctxUserEmail, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


// // Start of checking if given user is a SuperUser
// func(db *AuthMiddleWareStruct) HasSuperUserScope(next http.Handler) http.Handler{
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		jwttoken := r.Header.Get("Authorization")
// 		if jwttoken == ""{
// 			fmt.Println("No Authorization")
// 			return
// 		}
// 		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
// 		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("Testing key"), nil
// 		})
// 		if err != nil{
// 			fmt.Println("HasSuperUserScope failed")
// 			fmt.Println(err)
// 			helpers.ErrorJSON(w,err, 400)
// 			return
// 		}
// 		claims := token.Claims.(jwt.MapClaims)
// 		if claims["admin"] != "True"{
// 			err := errors.New("failed superUser check")
// 			helpers.ErrorJSON(w,err,400)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
// 		var exists bool
// 		db.db.QueryRow("SELECT UserID FROM tblAdminUsers WHERE UserID = ? AND SuperUser = 1", claims["userId"]).Scan(&exists)
// 		if exists == false{
// 			fmt.Println("User not in Admin, HasAdminScope has failed")
// 			err := errors.New("failed admin check")
// 			helpers.ErrorJSON(w,err, 400)
// 			return
// 		}
// 		next.ServeHTTP(w,r.WithContext(ctx))

// 	})
// }




// func(db *AuthMiddleWareStruct) HasAdminScope(next http.Handler) http.Handler{
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		jwttoken := r.Header.Get("Authorization")
// 		if jwttoken == ""{
// 			fmt.Println("No Authorization")
// 			return
// 		}
// 		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
// 		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("Testing key"), nil
// 		})
// 		if err != nil{
// 			fmt.Println("HasSuperUserScope failed")
// 			fmt.Println(err)
// 			helpers.ErrorJSON(w,err, 400)
// 			return
// 		}
// 		claims := token.Claims.(jwt.MapClaims)
// 		if claims["admin"] != "True"{
// 			err := errors.New("Failed Admin Check")
// 			helpers.ErrorJSON(w,err,400)
// 			return
// 		}
// 		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
// 		var exists bool
// 		db.db.QueryRow("SELECT UserID FROM tblAdminUsers WHERE UserID = ?", claims["userId"]).Scan(&exists)
// 		if exists == false{
// 			fmt.Println("User not in Admin, HasAdminScope has failed")
// 			err := errors.New("failed admin check")
// 			helpers.ErrorJSON(w,err, 400)
// 			return
// 		}
// 		next.ServeHTTP(w,r.WithContext(ctx))
// 	})
// }


















































// CustomClaims contains custom data we want from the token.
// type CustomClaims struct {
// 	Scope string `json:"scope"`
// }

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
// func (c CustomClaims) Validate(ctx context.Context) error {
// 	return nil
// }

// EnsureValidToken is a middleware that will check the validity of our JWT.
// The return is a function that also returns an http.Handler.
// func EnsureValidToken() func(next http.Handler) http.Handler {
// 	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
// 	if err != nil {
// 		log.Fatalf("Failed to parse the issuer url: %v", err)
// 	}

// 	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

// 	jwtValidator, err := validator.New(
// 		provider.KeyFunc,
// 		validator.RS256,
// 		issuerURL.String(),
// 		[]string{os.Getenv("AUTH0_API_AUDIENCE")},

// 		validator.WithCustomClaims(
// 			func() validator.CustomClaims {
// 				return &CustomClaims{}
// 			},
// 		),
// 		validator.WithAllowedClockSkew(time.Minute),
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to set up the jwt validator")
// 	}

// 	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
// 		log.Printf("Encountered error while validating JWT: %v", err)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
// 	}
// 	// After creating the errorHandler, we are going to be handling any possible issue that could arise from checking jwt. This is then passed on line 71.
// 	middleware := jwtmiddleware.New(
// 		jwtValidator.ValidateToken,
// 		jwtmiddleware.WithErrorHandler(errorHandler),
// 	)

// 	return func(next http.Handler) http.Handler {
// 		return middleware.CheckJWT(next)
// 	}
// }