package types

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	PasswordHashed string           `json:"-"`
	Subscription   SubscriptionType `json:"subscription"`
	RegisterDate   time.Time        `json:"register_date"`
	LastLogin      time.Time        `json:"last_login"`
	IsAdmin        bool             `json:"is_admin"`
	IsPaid         bool             `json:"is_paid"`
	Portfolio      Portfolio        `json:"portfolio"`
}

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func NewUser(params RegisterUser) (*User, error) {
	hashedPswrd, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:           "", // Can update name later
		Email:          params.Email,
		PasswordHashed: string(hashedPswrd),
		Subscription:   "nosubs",
		RegisterDate:   time.Now().Local(),
		LastLogin:      time.Now().Local(),
		IsAdmin:        false,
		IsPaid:         false,
	}, nil

}

func ValidatePassword(hashPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)) == nil
}

func CreateToken(user User) (string, error) {
	now := time.Now()
	validUntil := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"subscription":  user.Subscription,
		"register_date": user.RegisterDate,
		"last_login":    user.LastLogin,
		"is_admin":      user.IsAdmin,
		"is_paid":       user.IsPaid,
		"exp":           validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	tokenStr, err := token.SignedString([]byte(secret))

	if err != nil {

		return "", fmt.Errorf("failed to sign token %w", err)
	}

	return tokenStr, nil
}
