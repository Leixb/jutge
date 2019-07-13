package commands

import (
	"github.com/leixb/jutge/auth"

	"github.com/imroc/req"
)

type config struct {
	workDir string
	regex   string

	a           *auth.Credentials
	concurrency uint
}

var conf config

func init() {
	conf.workDir = "JutgeProblems"
	conf.regex = `[PGQX]\d{5}_(ca|en|es|fr|de)`
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

func SetRegex(regex string) {
	conf.regex = regex
}

func Concurrency() *uint {
	return &conf.concurrency
}
func Regex() *string {
	return &conf.regex
}

func getToken() (string, error) {
	if conf.a == nil {
		conf.a = auth.GetInstance()
	}
	return conf.a.TokenUID, nil
}

func getReq() (*req.Req, error) {
	if conf.a == nil {
		conf.a = auth.GetInstance()
	}
	return conf.a.R, nil
}
