package routes

import (
	"net/http"
	"os"

	"github.com/danilopolani/gocialite/structs"
	"github.com/gin-gonic/gin"
	"github.com/nsukmana-dev/restapi/config"
	"github.com/nsukmana-dev/restapi/models"
)

// func main() {
// 	router := gin.Default()

// 	router.GET("/", indexHandler)
// 	router.GET("/auth/:provider", RedirectHandler)
// 	router.GET("/auth/:provider/callback", CallbackHandler)

// 	router.Run("127.0.0.1:9090")
// }

// Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {

	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		// "linkedin": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/linkedin/callback",
		// },
		// "facebook": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/facebook/callback",
		// },
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_G"),
			"clientSecret": os.Getenv("CLIENT_SECRET_G"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
		// "bitbucket": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/bitbucket/callback",
		// },
		// "amazon": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/amazon/callback",
		// },
		// "slack": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/slack/callback",
		// },
		// "asana": {
		// 	"clientID":     "xxxxxxxxxxxxxx",
		// 	"clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		// 	"redirectURL":  "http://localhost:9090/auth/asana/callback",
		// },
	}

	providerScopes := map[string][]string{
		"github":    []string{"public_repo"},
		"linkedin":  []string{},
		"facebook":  []string{},
		"google":    []string{},
		"bitbucket": []string{},
		"amazon":    []string{},
		"slack":     []string{},
		"asana":     []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, token, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, user)

	c.JSON(200, gin.H{
		"data":    newUser,
		"token":   token,
		"message": "Berhasil login",
	})
}

func getOrRegisterUser(provider string, user *structs.User) models.User {

	var userData models.User

	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User{
			FullName: user.FullName,
			Email:    user.Email,
			SocialId: user.ID,
			Provider: provider,
			Avatar:   user.Avatar,
		}

		config.DB.Create(&newUser)
		return newUser
	} else {
		return userData
	}

}
