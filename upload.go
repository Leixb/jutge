package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Upload settings
type Upload struct {
	files       []string
	code        string
	compiler    string
	annotation  string
	check       bool
	concurrency int

	codes map[string]bool
}

// NewUpload return Upload object
func NewUpload() *Upload {
	return &Upload{codes: make(map[string]bool), concurrency: 3, compiler: "G++11"}
}

var compilers = []string{
	"BEEF", "Chicken", "CLISP", "Erlang", "F2C", "FBC", "FPC", "G++",
	"G++11", "GCC", "GCJ", "GDC", "GFortran", "GHC", "GNAT", "Go", "GObjC",
	"GPC", "Guile", "IVL08", "JDK", "Lua", "MakePRO2", "MonoCS", "nodejs",
	"P1++", "P2C", "Perl", "PHP", "PRO2", "Python", "Python3", "Quiz", "R",
	"Ruby", "RunHaskell", "RunPython", "Stalin", "Verilog", "WS",
}

// ConfigCommand configure kingpin options
func (u *Upload) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("upload", "Upload file to jutge.org").Action(u.Run)

	// Arguments
	cmd.Arg("file", "File to upload").Required().ExistingFilesVar(&u.files)

	// Flags
	cmd.Flag("compiler", "Compiler to use").Short('C').Default("G++11").EnumVar(&u.compiler, compilers...)
	cmd.Flag("code", "Problem code").Short('c').StringVar(&u.code)
	cmd.Flag("annotation", "Annotation").Short('a').Default("Uploaded with jutge_cli go").StringVar(&u.annotation)
	cmd.Flag("concurrency", "Number of simultaneous uploads").Default("3").IntVar(&u.concurrency)
	cmd.Flag("check", "Check veredict after upload").BoolVar(&u.check)
}

// Run the command
func (u *Upload) Run(*kingpin.ParseContext) error {
	err := u.UploadFiles()
	if err != nil {
		return err
	}
	if u.check {
		return u.CheckUploaded()
	}
	return nil
}

// CheckUploaded checks veredict of uploaded problems
func (u *Upload) CheckUploaded() error {
	var wg sync.WaitGroup
	sem := make(chan bool, u.concurrency)

	checker := NewCheck()

	for code, _ := range u.codes {
		sem <- true
		wg.Add(1)

		go func(c string) {
			defer func() { <-sem; wg.Done() }()
			for i := 0; i < 6; i++ {
				time.Sleep(time.Second * 5)
				veredict, err := checker.CheckLast(c)
				if err != nil {
					fmt.Println("Error checking", c, err)
					return
				}
				if veredict != "Not found" {
					fmt.Println(c, veredict)
					return
				}
			}
			fmt.Println(c, "Timed out")
		}(code)
	}

	wg.Wait()
	return nil
}

// UploadFiles upload all files in u.files
func (u *Upload) UploadFiles() error {
	var err error

	extractCode := u.code == ""

	var wg sync.WaitGroup
	sem := make(chan bool, u.concurrency)

	for _, file := range u.files {
		sem <- true
		wg.Add(1)
		go func(f string) {
			defer func() { <-sem; wg.Done() }()

			fmt.Println("Uploading:", f)

			code := u.code
			if extractCode {
				code, err = getCode(f)
				if err != nil {
					fmt.Println("Can't get code for file:", f)
					return
				}

			}

			err = u.UploadFile(f, code)
			if err != nil {
				fmt.Println("Upload failed", f, err)
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
func (u Upload) UploadFile(fileName, code string) error {
	token, err := Conf.getToken()
	if err != nil {
		return err
	}

	param := req.Param{
		"annotation":  u.annotation,
		"compiler_id": u.compiler,
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

	rq, err := Conf.getReq()
	if err != nil {
		return err
	}

	_, err = rq.Post("https://jutge.org/problems/"+code+"/submissions", param, file)
	return err
}
