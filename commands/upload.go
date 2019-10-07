package commands

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/imroc/req"
)

type upload struct {
	Code       string
	Compiler   string
	Annotation string
	Check      bool

	codes map[string]bool
}

// NewUpload return upload object
func NewUpload() *upload {
	return &upload{codes: make(map[string]bool), Compiler: "G++11"}
}

// GetCompilers list with all valid compilers for upload
func GetCompilers() []string {
	return compilers
}

var compilers = []string{
	"BEEF", "Chicken", "CLISP", "Erlang", "F2C", "FBC", "FPC", "G++",
	"G++11", "GCC", "GCJ", "GDC", "GFortran", "GHC", "GNAT", "Go", "GObjC",
	"GPC", "Guile", "IVL08", "JDK", "Lua", "MakePRO2", "MonoCS", "nodejs",
	"P1++", "P2C", "Perl", "PHP", "PRO2", "Python", "Python3", "Quiz", "R",
	"Ruby", "RunHaskell", "RunPython", "Stalin", "Verilog", "WS",
}

// UploadFiles upload all files in files
func (u *upload) UploadFiles(files []string) error {
	var err error

	extractCode := u.Code == ""

	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	for _, file := range files {
		sem <- true
		wg.Add(1)
		go func(f string) {
			defer func() { <-sem; wg.Done() }()

			fmt.Println(" - Uploading:", f)

			code := u.Code
			if extractCode {
				code, err = getCode(f)
				if err != nil {
					fmt.Println(" ! Can't get code for file:", f)
					return
				}

			}

			err = u.UploadFile(f, code)
			if err != nil {
				fmt.Println(" ! Upload failed", f, err)
				return
			}
			// Add code to set so it can be checked later
			u.codes[code] = true

		}(file)
	}

	wg.Wait()

	return err
}

// UploadFile submit file to jutge.org
func (u *upload) UploadFile(fileName, code string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	param := req.Param{
		"annotation":  u.Annotation,
		"compiler_id": u.Compiler,
		"submit":      "submit",
		"token_uid":   token,
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file := req.FileUpload{
		File:      f,
		FieldName: "file",
		FileName:  "program",
	}

	rq, err := getReq()
	if err != nil {
		return err
	}

	_, err = rq.Post("https://jutge.org/problems/"+code+"/submissions", param, file)
	return err
}

// CheckUploaded checks veredict of uploaded problems
func (u *upload) CheckUploaded() error {
	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	checker := NewCheck()

	for code := range u.codes {
		sem <- true
		wg.Add(1)

		go func(c string) {
			defer func() { <-sem; wg.Done() }()
			for i := 0; i < 6; i++ {
				time.Sleep(time.Second * 5)
				veredict, err := checker.CheckLast(c)
				if err != nil {
					fmt.Println(" ! Error checking", c, err)
					return
				}
				if veredict != "Not found" {
					fmt.Println(" -", c, veredict)
					return
				}
			}
			fmt.Println(" !", c, "Timed out")
		}(code)
	}

	wg.Wait()
	return nil
}
