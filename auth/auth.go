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

	"github.com/PuerkitoBio/goquery"
	"github.com/howeyc/gopass"
	"github.com/imroc/req"
)

type Auth struct {
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

func Login(credentials ...string) (Auth, error) {
	var ldata loginData
	switch len(credentials) {
	case 0:
	case 1:
		ldata.Email = credentials[0]
	case 2:
		ldata.Email = credentials[0]
		ldata.Password = credentials[1]
	default:
		return Auth{}, errors.New("Too many parameters")
	}

	var a Auth

	if a.loadTmp() {
		return a, nil
	}

	err := ldata.promptMissing()
	if err != nil {
		return a, err
	}

	a, err = login(ldata)
	if err != nil {
		return a, err
	}

	err = a.saveTmp()

	return a, err
}

func login(ldata loginData) (auth Auth, err error) {

	auth.Username = ldata.Email

	auth.R = req.New()

	auth.R.EnableCookie(true)

	params := req.Param{
		"email":    ldata.Email,
		"password": ldata.Password,
		"submit":   "submit",
	}

	resp, err := auth.R.Post("https://jutge.org", params)
	if err != nil {
		return
	}

	if len(resp.Response().Header.Get("Set-Cookie")) > 0 {
		err = errors.New("Failed Login")
		return
	}

	err = auth.setTokenUID()
	if err != nil {
		return
	}

	return
}

func (a *Auth) setTokenUID() error {

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
		return errors.New("Token UID not found")
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
func (a *Auth) saveTmp() error {

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
func (a *Auth) loadTmp() bool {
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
		fmt.Println("Cookie Expired")
		return false
	}

	return true
}
