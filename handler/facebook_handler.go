package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/anime454/facebook-login/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type facebookHandler struct {
	fbHandler service.FaceBookService
}

func NewFacebookHandler(fb service.FaceBookService) facebookHandler {
	return facebookHandler{fbHandler: fb}
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

func (fbHandler facebookHandler) Callback() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		state, _ := c.GetQuery("state")
		fmt.Println("DEBUG: state > ", state)

		code, _ := c.GetQuery("code")
		fmt.Println("DEBUG: code > ", code)

		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		resp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" +
			url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Printf("Get: %s\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ReadAll: %s\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		log.Printf("parseResponseBody: %s\n", string(response))

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return fn
}

func (fb facebookHandler) LoginPage() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		const htmlIndex = `<html><body> Logged in with <a href="/login">facebook</a> </body></html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.Write([]byte(htmlIndex))
	}
	return fn
}

func (fb facebookHandler) Login() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
		if err != nil {
			log.Fatal("Parse: ", err)
		}
		parameters := url.Values{}
		parameters.Add("client_id", oauthConf.ClientID)
		parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
		parameters.Add("redirect_uri", oauthConf.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", oauthStateString)
		Url.RawQuery = parameters.Encode()
		url := Url.String()
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
	return fn
}
