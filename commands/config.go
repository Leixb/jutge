package commands

import (
	"regexp"

	"github.com/leixb/jutge/auth"

	"github.com/imroc/req"
)

type config struct {
	workDir string
	regex   *regexp.Regexp

	username, password string

	a           *auth.Credentials
	concurrency uint
}

var conf config

func init() {
	conf.workDir = "JutgeProblems"
	conf.regex = regexp.MustCompile(`[PGQX]\d{5}_(ca|en|es|fr|de)`)
	conf.concurrency = 3
}

func SetWorkDir(dir string) {
	conf.workDir = dir
}

func SetConcurrency(n uint) {
	conf.concurrency = n
	if n == 0 {
		conf.concurrency = 50
	}
}

func SetRegex(regex string) (err error) {
	conf.regex, err = regexp.Compile(regex)
	return
}

func Concurrency() *uint {
	return &conf.concurrency
}

func Regex() **regexp.Regexp {
	return &conf.regex
}

func WorkDir() *string {
	return &conf.workDir
}

func SetUsername(username string) {
	conf.username = username
}
func SetPassword(password string) {
	conf.password = password
}

func getToken() (string, error) {
	if conf.a == nil {
		conf.a = auth.GetInstance(conf.username, conf.password)
	}
	return conf.a.TokenUID, nil
}

func getReq() (*req.Req, error) {
	if conf.a == nil {
		conf.a = auth.GetInstance(conf.username, conf.password)
	}
	return conf.a.R, nil
}
