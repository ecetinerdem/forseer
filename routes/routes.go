package routes

import (
	"fmt"
	"net/http"
)

func HandleGreeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handle users")
}
