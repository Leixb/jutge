package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Upload settings
type Upload struct {
	compiler    string
	files       []string
	code        string
	annotation  string
	concurrency int
}

// NewUpload return Upload object
func NewUpload() *Upload {
	return &Upload{}
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
}

// Run the command
func (u *Upload) Run(*kingpin.ParseContext) error {
	return u.UploadFiles()
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
			if extractCode {
				u.code, err = getCode(f)
				if err != nil {
					fmt.Println("Can't get code for file:", f)
					return
				}
			}
			err = u.UploadFile(f)
			if err != nil {
				fmt.Println("Upload failed", f, err)
			}
		}(file)
	}

	wg.Wait()

	return err
}

// UploadFile submit file to jutge.org
func (u Upload) UploadFile(fileName string) error {
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

	_, err = rq.Post("https://jutge.org/problems/"+u.code+"/submissions", param, file)
	return err
}
