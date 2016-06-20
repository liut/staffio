package exmail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"lcgc/platform/staffio/models"
)

var (
	errEmptyAuths = errors.New("empty auths")
	api_auths     string
)

const (
	url_token    = "https://exmail.qq.com/cgi-bin/token"
	url_user_get = "http://openapi.exmail.qq.com:12211/openapi/user/get"
	url_newcount = "http://openapi.exmail.qq.com:12211/openapi/mail/newcount"
)

func init() {
	api_auths = os.Getenv("EXMAIL_API_AUTHS")
}

type OpenType uint8

const (
	OTIgnore   OpenType = 0
	OTEnabled  OpenType = 1
	OTDisabled OpenType = 2
)

func GetStaff(email string) (*models.Staff, error) {
	user, err := requestUserGet(email)
	if err != nil {
		return nil, err
	}

	sn, gn := models.SplitName(user.Name)
	log.Printf("%q %q %q", user.Name, sn, gn)

	return &models.Staff{
		Uid:            strings.Split(user.Alias, "@")[0],
		Email:          user.Alias,
		CommonName:     user.Name,
		Surname:        sn,
		GivenName:      gn,
		EmployeeNumber: user.ExtId,
		EmployeeType:   user.Title,
		Mobile:         user.Mobile,
		Gender:         user.Gender,
	}, nil
}

/*{
"Alias": " test2@gzservice.com",
"Name": "鲍勃",
"Gender": 1,
"SlaveList": "bb@gzdev.com,bo@gzdev.com",
"Position": "工程师",
"Tel": "62394",
"Mobile": "",
"ExtId": "100",
"PartyList": {
	"Count": 3,
	"List": [{ "Value":"部门 a" }
		,{ "Value":"部门 B/部门 b" }
		,{"Value":"部门 c" }
}}*/
type userResp struct {
	Alias    string        `json:"Alias"`
	Name     string        `json:"Name"`
	Aliases  string        `json:"SlaveList"`
	Gender   models.Gender `json:"Gender"`
	Title    string        `json:"Position"`
	ExtId    string        `json:"ExtId"`
	Tel      string        `json:"Tel"`
	Mobile   string        `json:"Mobile"`
	OpenType OpenType      `json:"OpenType"`
}

type apiError struct {
	Arg     string `json:"arg"`
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"error"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s:%s %q", e.ErrCode, e.ErrMsg, e.Arg)
}

type newCount struct {
	Alias    string
	NewCount json.Number
}

func RequestMailNewCount(alias string) (int, error) {
	token, err := requestAccessToken()
	if err != nil {
		return 0, err
	}
	auths := "Bearer " + token
	resp, err := doHTTP("POST", url_newcount, auths, bytes.NewBufferString("alias="+alias))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return 0, err
	}

	log.Printf("resp: %s", resp)

	obj := &newCount{}

	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return 0, e
	}

	count, err := obj.NewCount.Int64()
	if err != nil {
		log.Print(err)
	}
	return int(count), nil
}

func requestUserGet(alias string) (*userResp, error) {
	token, err := requestAccessToken()
	if err != nil {
		return nil, err
	}
	auths := "Bearer " + token
	resp, err := doHTTP("POST", url_user_get, auths, bytes.NewBufferString("alias="+alias))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return nil, err
	}

	log.Printf("resp: %s", resp)

	exErr := &apiError{Arg: alias}
	if e := json.Unmarshal(resp, exErr); e != nil {
		log.Printf("unmarshal api err %s", e)
		return nil, e
	}

	if exErr.ErrCode != "" {
		log.Printf("apiError %s", exErr)
		return nil, exErr
	}

	obj := &userResp{}
	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return nil, e
	}

	return obj, nil
}

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func requestAccessToken() (token string, err error) {
	if api_auths == "" {
		err = errEmptyAuths
		log.Print(err)
		return
	}
	auths := "Basic " + api_auths
	// log.Printf("auths: %s", auths)
	body_str := "grant_type=client_credentials"
	resp, err := doHTTP("POST", url_token, auths, bytes.NewBufferString(body_str))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return
	}

	// log.Printf("%s", resp)
	obj := &tokenResp{}
	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal err %s", e)
		return
	}

	token = obj.AccessToken
	return
}

func doHTTP(method, url string, auths string, body io.Reader) ([]byte, error) {

	// log.Printf("doHTTP: %s %s", method, url)

	req, e := http.NewRequest(method, url, body)
	if e != nil {
		log.Println(e, method, url)
		return nil, e
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if auths != "" {
		req.Header.Set("Authorization", auths)
	}

	c := &http.Client{}
	resp, e := c.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		log.Printf("http code error %d, %s", resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	return rbody, nil
}
