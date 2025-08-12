package types

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	Password     string
	IsAdmin      bool
	IsPaid       bool
	Subscription SubscriptionType
	RegisterDate time.Time
	LastLogin    time.Time
	//TODO: Add Portfolio slice and radar slice which we will take from a stocks api
	//have to create types according to fetched data
}

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"passoword"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"passoword"`
}

type LoginUserResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
