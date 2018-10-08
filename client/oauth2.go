package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

const (
	cKeyState = "staffio_state"
	cKeyToken = "staffio_token"
	cKeyUser  = "staffio_user"
)

var (
	conf           *oauth2.Config
	oAuth2Endpoint oauth2.Endpoint
	infoUrl        string
)

func init() {
	prefix := envOr("STAFFIO_PREFIX", "https://staffio.work")
	oAuth2Endpoint = oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("%s/%s", prefix, "authorize"),
		TokenURL: fmt.Sprintf("%s/%s", prefix, "token"),
	}
	clientID := envOr("STAFFIO_CLIENT_ID", "")
	clientSecret := envOr("STAFFIO_CLIENT_SECRET", "")
	if clientID == "" || clientSecret == "" {
		log.Print("Warning: STAFFIO_CLIENT_ID or STAFFIO_CLIENT_SECRET not found in environment")
	}
	infoUrl = fmt.Sprintf("%s/%s", prefix, "info/me")
	redirectURL := envOr("STAFFIO_REDIRECT_URL", "")
	scopes := strings.Split(envOr("STAFFIO_SCOPES", ""), ",")
	if clientID != "" && clientSecret != "" {
		Setup(redirectURL, clientID, clientSecret, scopes)
	}
}

func randToken() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GetOAuth2Config() *oauth2.Config {
	return conf
}

// Setup oauth2 config
func Setup(redirectURL, clientID, clientSecret string, scopes []string) {
	conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     oAuth2Endpoint,
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fakeUser := &User{Name: randToken()}
	fakeUser.Refresh()
	state, _ := fakeUser.Encode()
	http.SetCookie(w, &http.Cookie{
		Name:     cKeyState,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
	})
	// g6F1oKFu2SxkalR3cms5dkIxSW1KcmppWUFxOThnNUlBS0lFNE5EV1VpUkgzMVdxc1VBPaFo0luyTHg
	// g6F1oKFusDZVa0J6MDNVLTZ1Q2lKN1ShaNJbsldv
	sess := SessionLoad(r)
	sess.Set(cKeyState, state)
	SessionSave(sess, w)
	location := GetAuthCodeURL(state)
	w.Header().Set("refresh", fmt.Sprintf("1; %s", location))
	w.Write([]byte("<html><title>Staffio</title> <body style='padding: 2em;'> <p>Waiting...</p> <a href='" +
		location + "'><button style='font-size: 14px;'> Login with Staffio! </button></a></body></html>"))
}

func GetAuthCodeURL(state string) string {
	return conf.AuthCodeURL(state)
}

func envOr(key, dft string) string {
	v := os.Getenv(key)
	if v == "" {
		return dft
	}
	return v
}
