package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ecetinerdem/forseer/types"
	"github.com/golang-jwt/jwt/v5"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *types.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *types.User {
	val := ctx.Value(userKey)

	user, ok := val.(*types.User)

	if !ok {
		return nil
	}
	return user
}

func UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 || headerParts[0] == "Bearer" {
			http.Error(w, "Unauthorized: invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		claims, err := ParseToken(token)

		if err != nil {
			http.Error(w, "Unauthorized: invalid token header", http.StatusUnauthorized)
			return
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			http.Error(w, "Unauthorized: token expired", http.StatusUnauthorized)
			return
		}

		var user types.User

		user.Email = claims["email"].(string)
		user.ID = claims["id"].(string)

		ctx := r.Context()

		ctx = WithUser(ctx, &user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized: token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized: invalid claims format")
	}

	return claims, nil
}
