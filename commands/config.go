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

// SetWorkDir sets the location where to download and check for problem files and the database
func SetWorkDir(dir string) {
	conf.workDir = dir
}

// SetConcurrency sets the number of concurrent goroutines that can be run at the same time
func SetConcurrency(n uint) {
	conf.concurrency = n
	if n == 0 {
		conf.concurrency = 50
	}
}

// Setregex sets the regex used to validate and extract problem codes from filenames
func SetRegex(regex string) (err error) {
	conf.regex, err = regexp.Compile(regex)
	return
}

// Concurrency returns reference to conf.concurrency. Use with caution.
func Concurrency() *uint {
	return &conf.concurrency
}

// Regex returns reference to conf.regex. Use with caution.
func Regex() **regexp.Regexp {
	return &conf.regex
}

// WorkDir returns reference to conf.WorkDir. Use with caution.
func WorkDir() *string {
	return &conf.workDir
}

// SetUsername sets username for login
func SetUsername(username string) {
	conf.username = username
}

// SetPassword sets password for login
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
