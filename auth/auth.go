package auth

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

// TokenNotFound error returned when the parser can't find the TokenUID, this usually means
// that the login failed or cookies expired
type TokenNotFound struct{}

func (*TokenNotFound) Error() string { return "Token uid not found" }

var (
	singleton *Credentials
	once      sync.Once
)

// Credentials object that holds the necessary information to make a request.
// It does not save the user password, only a tooken and cookies.
type Credentials struct {
	Username string

	// TokenUID is needed when uploading files
	TokenUID string

	// R is a imroc/req.Req that allows to perform requests with the proper cookies set.
	R *req.Req
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

// GetInstance creates and returns a new Credentials instance.
// It accepts at most 2 parameters: email and password. If any of them are omitted
// the program will prompt the user to input them through stdout.
// Note that you can not omit the email and not the password.
// It panics when login fails.
func GetInstance(creds ...string) *Credentials {
	once.Do(func() {
		if len(creds) > 2 {
			panic("Too many argumetns")
		}

		var ldata loginData

		if len(creds) > 0 {
			ldata.Email = creds[0]
		}
		if len(creds) > 1 {
			ldata.Password = creds[1]
		}

		var err error
		singleton, err = newInstance(ldata)
		if err != nil {
			panic("Failed Login: " + err.Error())
		}
	})
	return singleton
}

func newInstance(ldata loginData) (*Credentials, error) {

	cred := &Credentials{TokenUID: "", Username: ""}
	if cred.loadTmp() {
		return cred, nil
	}

	err := ldata.promptMissing()
	if err != nil {
		return nil, err
	}
	cred, err = login(ldata)
	if err != nil {
		return nil, err
	}

	err = cred.saveTmp()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Save credentials error:", err)
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

func (a *Credentials) getTokenUID() (string, bool, error) {
	resp, err := a.R.Get("https://jutge.org/problems/P68688_ca")
	if err != nil {
		return "", false, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp.Response())
	if err != nil {
		return "", false, err
	}

	value, ok := doc.Find(".col-sm-4 > input:nth-child(1)").Attr("value")
	return value, ok, nil
}

func (a *Credentials) setTokenUID() error {
	var ok bool
	var err error

	a.TokenUID, ok, err = a.getTokenUID()
	if err != nil {
		return err
	}

	if !ok {
		return &TokenNotFound{}
	}

	return nil
}

func (a *Credentials) IsExpired() bool {
	_, found, _ := a.getTokenUID()
	return !found
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
		fmt.Fprintln(os.Stderr, "open", err)
		return false
	}
	defer file.Close()

	var save persist

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "readall", err)
		return false
	}

	err = json.Unmarshal(data, &save)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unmarshall", err)
		return false
	}

	u, err := url.Parse("https://jutge.org")
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse", err)
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
			fmt.Fprintln(os.Stderr, "Cookie Expired")
		default:
			fmt.Fprintln(os.Stderr, "Connection error")
		}
		return false
	}

	return true
}
