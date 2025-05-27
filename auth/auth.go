package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"graphQlDemo/ent"
	"graphQlDemo/ent/user"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

var secretKey = []byte("my-super-secret-key")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["user"].(string), nil
	}

	return "", errors.New("token inválido")
}

func UserFromContext(ctx context.Context) (*ent.User, bool) {
	user, ok := ctx.Value(userContextKey).(*ent.User)
	return user, ok
}

func isPublicOperation(opName, query string) bool {
	publicOps := []string{"login", "register"}
	opLower := strings.ToLower(opName)
	queryLower := strings.ToLower(query)

	for _, op := range publicOps {
		if strings.Contains(opLower, op) || strings.Contains(queryLower, op) {
			return true
		}
	}
	return false
}

func Middleware(client *ent.Client, next *handler.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var reqBody struct {
				OperationName string `json:"operationName"`
				Query         string `json:"query"`
			}

			if json.Unmarshal(bodyBytes, &reqBody) == nil {
				if isPublicOperation(reqBody.OperationName, reqBody.Query) {
					next.ServeHTTP(w, r)
					return
				}
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		username, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		user, err := client.User.
			Query().
			Where(user.Username(username)).
			Only(r.Context())
		if err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
