package commands

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/Leixb/jutge/auth"
	"github.com/imroc/req"
)

type jutge struct {
	url   string
	auth  *auth.Credentials
	regex *regexp.Regexp
	folder string
}

type JutgeConfig struct {
	URL   string
	Regex *regexp.Regexp
	Auth  *auth.Credentials
	Folder string
}

func Jutge(config *JutgeConfig) *jutge {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	// Default values
	j := &jutge{
		url:   "https://jutge.org/",
		auth:  nil,
		regex: regexp.MustCompile(`[PGQX]\d{5}_(ca|en|es|fr|de)`),
		folder: filepath.Join(home, "jutge"),
	}

	if config == nil {
		return j
	}

	if config.Regex != nil {
		j.regex = config.Regex
	}

	if config.URL != "" {
		j.url = config.URL
	}

	if config.Auth != nil {
		j.auth = config.Auth
	}

	if config.Folder != "" {
		j.folder = config.Folder
	}

	return j
}

func (j *jutge) GetURL() string {
	return j.url
}

func (j *jutge) GetToken() string {
	if j.auth == nil {
		j.Login("", "")
	}
	return j.auth.TokenUID
}

func (j *jutge) GetReq() *req.Req {
	if j.auth == nil {
		j.Login("", "")
	}
	return j.auth.R
}

func (j *jutge) Login(username, password string) {
	if j.auth == nil || (username != "" && j.auth.Username != username) || j.auth.IsExpired() {
		hasUser := username != ""
		hasPass := password != ""

		if hasUser && hasPass {
			j.auth = auth.GetInstance(username, password)
		} else if hasUser {
			j.auth = auth.GetInstance(username)
		} else {
			j.auth = auth.GetInstance()
		}
	}
}
