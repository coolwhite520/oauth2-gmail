package api

import (
	"encoding/json"
	"log"
	"oauth2-gmail/database"
	"oauth2-gmail/model"
)

// InitializeProfile Initializes the user in the database
func InitializeProfile(accessToken string, refreshToken string) {
	userResponse, _ := CallAPIMethod("GET", model.ProfileEndpointRoot, "/userinfo", accessToken, "", nil, "")
	user := model.User{}
	user.AccessToken = accessToken
	user.AccessTokenActive = 1
	user.RefreshToken = refreshToken
	json.Unmarshal([]byte(userResponse), &user)
	user.UserPrincipalName = user.Mail
	database.InsertUser(user)
	log.Println(user.Mail, "success log in.")
}
