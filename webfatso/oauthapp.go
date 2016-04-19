package webfatso

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"

	. "lcgc/platform/staffio/settings"
)

// addr = localhost:3000
func AppDemo(router *mux.Router) {

	const (
		demoId     = "1234"
		demoSecret = "aabbccdd"
	)

	addr := Settings.HttpListen

	// Application home endpoint
	router.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>"))

		w.Write([]byte(fmt.Sprintf("<a href=\"/authorize?response_type=code&client_id=1234&state=xyz&scope=everything&redirect_uri=%s\">Code</a><br/>", url.QueryEscape("http://"+addr+"/appauth/code"))))
		w.Write([]byte(fmt.Sprintf("<a href=\"/authorize?response_type=token&client_id=1234&state=xyz&scope=everything&redirect_uri=%s\">Implict</a><br/>", url.QueryEscape("http://"+addr+"/appauth/token"))))
		w.Write([]byte(fmt.Sprintf("<a href=\"/appauth/password\">Password</a><br/>")))
		w.Write([]byte(fmt.Sprintf("<a href=\"/appauth/client_credentials\">Client Credentials</a><br/>")))
		w.Write([]byte(fmt.Sprintf("<a href=\"/appauth/assertion\">Assertion</a><br/>")))

		w.Write([]byte("</body></html>"))
	})

	// Application destination - CODE
	router.HandleFunc("/appauth/code", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		code := r.Form.Get("code")

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))
		defer w.Write([]byte("</body></html>"))

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&state=xyz&redirect_uri=%s&code=%s",
			url.QueryEscape("http://"+addr+"/appauth/code"), url.QueryEscape(code))

		// if parse, download and parse json
		if r.Form.Get("doparse") == "1" {
			err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
				&osin.BasicAuth{demoId, demoSecret}, jr)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.Write([]byte("<br/>"))
			}
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		// output links
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Goto Token URL</a><br/>", aurl)))

		cururl := *r.URL
		curq := cururl.Query()
		curq.Add("doparse", "1")
		cururl.RawQuery = curq.Encode()
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}

		if at, ok := jr["access_token"]; ok {
			rurl := fmt.Sprintf("/appauth/info?code=%s", at)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
		}
	})

	// Application destination - TOKEN
	router.HandleFunc("/appauth/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - TOKEN<br/>"))

		w.Write([]byte("Response data in fragment - not acessible via server - Nothing to do"))

		w.Write([]byte("</body></html>"))
	})

	// Application destination - PASSWORD
	router.HandleFunc("/appauth/password", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - PASSWORD<br/>"))

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/token?grant_type=password&scope=everything&username=%s&password=%s",
			"test", "test")

		// download token
		err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
			&osin.BasicAuth{Username: demoId, Password: demoSecret}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}

		if at, ok := jr["access_token"]; ok {
			rurl := fmt.Sprintf("/appauth/info?code=%s", at)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
		}

		w.Write([]byte("</body></html>"))
	})

	// Application destination - CLIENT_CREDENTIALS
	router.HandleFunc("/appauth/client_credentials", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CLIENT CREDENTIALS<br/>"))

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/token?grant_type=client_credentials")

		// download token
		err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
			&osin.BasicAuth{Username: demoId, Password: demoSecret}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}

		if at, ok := jr["access_token"]; ok {
			rurl := fmt.Sprintf("/appauth/info?code=%s", at)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
		}

		w.Write([]byte("</body></html>"))
	})

	// Application destination - ASSERTION
	router.HandleFunc("/appauth/assertion", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - ASSERTION<br/>"))

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/token?grant_type=assertion&assertion_type=urn:osin.example.complete&assertion=osin.data")

		// download token
		err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
			&osin.BasicAuth{Username: demoId, Password: demoSecret}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}

		if at, ok := jr["access_token"]; ok {
			rurl := fmt.Sprintf("/appauth/info?code=%s", at)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
		}

		w.Write([]byte("</body></html>"))
	})

	// Application destination - REFRESH
	router.HandleFunc("/appauth/refresh", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - REFRESH<br/>"))
		defer w.Write([]byte("</body></html>"))

		code := r.Form.Get("code")

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/token?grant_type=refresh_token&refresh_token=%s", url.QueryEscape(code))

		// download token
		err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
			&osin.BasicAuth{Username: demoId, Password: demoSecret}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}

		if at, ok := jr["access_token"]; ok {
			rurl := fmt.Sprintf("/appauth/info?code=%s", at)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
		}
	})

	// Application destination - INFO
	router.HandleFunc("/appauth/info", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - INFO<br/>"))
		defer w.Write([]byte("</body></html>"))

		code := r.Form.Get("code")

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf("/info/me?code=%s", url.QueryEscape(code))

		// download token
		err := DownloadAccessToken(fmt.Sprintf("http://"+addr+"%s", aurl),
			&osin.BasicAuth{Username: demoId, Password: demoSecret}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		if rt, ok := jr["refresh_token"]; ok {
			rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
			w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
		}
	})

}

func HandleLoginPage(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	if r.Method == "POST" && r.Form.Get("login") == "test" && r.Form.Get("password") == "test" {
		return true
	}

	w.Write([]byte("<html><body>"))

	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?response_type=%s&client_id=%s&state=%s&redirect_uri=%s\" method=\"POST\">",
		ar.Type, ar.Client.GetId(), ar.State, url.QueryEscape(ar.RedirectUri))))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return false
}

func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return fmt.Errorf("Invalid status code: %d", presp.StatusCode)
	}

	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
