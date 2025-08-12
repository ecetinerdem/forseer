package types

import "time"

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
	// TODO: Add Portfolio slice and radar slice which we will take from a stocks API
	// have to create types according to fetched data
}

type CreateUser struct {
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
