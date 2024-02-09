package database

import (
	"log"
	"os"

	"github.com/odunlamizo/ovalfi/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var users []model.User

func init() {
	// create default users all with password ovalfi
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("ovalfi"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	users = append(users, model.User{Username: "Kenny", Password: hashedPassword})
	users = append(users, model.User{Username: "Ttilaayo", Password: hashedPassword})
	users = append(users, model.User{Username: "Viktor", Password: hashedPassword})
	users = append(users, model.User{Username: "Shegzy", Password: hashedPassword})
	users = append(users, model.User{Username: "Mhiddey", Password: hashedPassword})
}

// returns user with the specified username
func GetUser(username string) model.User {
	var user model.User
	for _, u := range users {
		if u.Username == username {
			user = u
			break
		}
	}
	return user
}
