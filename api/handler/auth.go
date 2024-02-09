package handler

import (
	"encoding/json"
	"net/http"

	"github.com/odunlamizo/ovalfi/internal/database"
	"github.com/odunlamizo/ovalfi/internal/model"
	"github.com/odunlamizo/ovalfi/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// login handler, returns jwt token if credentials are correct
func Login(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")
	if username == "" || password == "" {
		errorMessage := model.ResponseMessage{Message: "Missing parameter(s)"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	var user model.User
	databaseUser := database.GetUser(username)
	user = databaseUser
	if user.Username == "" {
		errorMessage := model.ResponseMessage{Message: "Username not found"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	hashedPassword := user.Password
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Incorrect password"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
		return
	}
	s, err := util.CreateToken(username)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error creating token"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	token := model.Token{AccessToken: s}
	responseBody, err := json.Marshal(token)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
