package auth

// package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/howeyc/gopass"
	"github.com/imroc/req"
)

// TokenNotFound error
type TokenNotFound struct{}

func (*TokenNotFound) Error() string { return "Token uid not found" }

var (
	singleton *Credentials
	once      sync.Once
)

// Auth object
type Credentials struct {
	Username string

	TokenUID string
	R        *req.Req
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type persist struct {
	Phpsessid string `json:"PHPSESSID"`
	TokenUID  string `json:"token_uid"`
	Email     string `json:"email"`
}

// Create new Auth instance
func GetInstance() *Credentials {
	once.Do(func() {
		var err error
		singleton, err = newInstance()
		if err != nil {
			panic("Failed Login: " + err.Error())
		}
	})
	return singleton
}

func newInstance() (*Credentials, error) {
	cred := &Credentials{TokenUID: "", Username: ""}
	if cred.loadTmp() {
		return cred, nil
	}

	var ldata loginData
	err := ldata.promptMissing()
	if err != nil {
		return nil, err
	}
	cred, err = login(ldata)

	err = cred.saveTmp()
	if err != nil {
		fmt.Println("Save credentials error:", err)
	}
	return cred, nil
}

func login(ldata loginData) (cred *Credentials, err error) {

	cred = &Credentials{Username: ldata.Email, R: req.New()}

	cred.R.EnableCookie(true)

	params := req.Param{
		"email":    ldata.Email,
		"password": ldata.Password,
		"submit":   "submit",
	}

	resp, err := cred.R.Post("https://jutge.org", params)
	if err != nil {
		return
	}

	if len(resp.Response().Header.Get("Set-Cookie")) > 0 {
		err = errors.New("Invalid Username / Password")
		return
	}

	err = cred.setTokenUID()
	if err != nil {
		return
	}

	return
}

func (a *Credentials) setTokenUID() error {

	resp, err := a.R.Get("https://jutge.org/problems/P68688_ca")
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(resp.Response())
	if err != nil {
		return err
	}

	var ok bool
	a.TokenUID, ok = doc.Find(".col-sm-4 > input:nth-child(1)").Attr("value")

	if !ok {
		return &TokenNotFound{}
	}

	return nil
}

func (ldata *loginData) promptMissing() error {
	if ldata.Email == "" {
		fmt.Print("Email: ")
		_, err := fmt.Scan(&ldata.Email)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Email:", ldata.Email)
	}

	if ldata.Password == "" {
		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		ldata.Password = string(pass)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Password:", ldata.Password)
	}
	return nil
}

func tmpFilename() string {
	return os.TempDir() + "/jutge_go_auth"
}

// Save credentials
func (a *Credentials) saveTmp() error {

	u, err := url.Parse("https://jutge.org")
	if err != nil {
		return err
	}

	cookies := a.R.Client().Jar.Cookies(u)

	save := persist{TokenUID: a.TokenUID, Email: a.Username}

	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			save.Phpsessid = cookie.Value
			break
		}
	}

	file, err := os.Create(tmpFilename())
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(save)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

// Restore saved credentials
func (a *Credentials) loadTmp() bool {
	file, err := os.Open(tmpFilename())
	if err != nil {
		fmt.Println("open", err)
		return false
	}
	defer file.Close()

	var save persist

	data, err := ioutil.ReadAll(file)

	err = json.Unmarshal(data, &save)
	if err != nil {
		fmt.Println("unmarshall", err)
		return false
	}

	u, err := url.Parse("https://jutge.org")
	if err != nil {
		fmt.Println("parse", err)
		return false
	}

	a.TokenUID = save.TokenUID
	a.Username = save.Email

	a.R = req.New()
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{Name: "PHPSESSID", Value: save.Phpsessid})
	a.R.Client().Jar.SetCookies(u, cookies)

	err = a.setTokenUID()
	if err != nil {
		switch err.(type) {
		case *TokenNotFound:
			fmt.Println("Cookie Expired")
		default:
			fmt.Println("Connection error")
		}
		return false
	}

	return true
}
