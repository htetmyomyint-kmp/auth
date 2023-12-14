package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/htetmyomyint-kmp/leave-tracker/auth/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
	state       string
)

func init() {
	oauthConfig = &oauth2.Config{
		ClientID:     "572115871227-cqu5vtr316jt65fdtsmrfvgse2lk010h.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/api/auth", // Update with your redirect URL
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

}

type GoogleAuthHandler struct {
	Logger   *log.Logger
	DBClient data.UserDatabase
}

func NewGoogleAuthHandler(l *log.Logger, db data.UserDatabase) GoogleAuthHandler {
	return GoogleAuthHandler{
		Logger:   l,
		DBClient: db,
	}
}

func (g *GoogleAuthHandler) Signup(c *gin.Context) {
	state = "sign-up"
	url := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *GoogleAuthHandler) Login(c *gin.Context) {
	state = "login"
	url := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *GoogleAuthHandler) Auth(c *gin.Context) {
	receivedState := c.Query("state")
	log.Println("state ", receivedState)
	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	userInfo, err := getUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	if receivedState == "sign-up" {
		if _, err := g.DBClient.CreateUser(data.User{
			Email:      userInfo["email"].(string),
			Name:       userInfo["name"].(string),
			ProfileURL: userInfo["picture"].(string),
		}); err != nil {
			c.JSON(http.StatusInternalServerError, "create error")
			return
		}
	}

	c.JSON(http.StatusOK, userInfo)

}

func getUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := oauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
