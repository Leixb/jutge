package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge_go/auth"
)

// Upload object that wraps its settings
type Upload struct {
	compiler string
	files    []string
	code     string
}

// NewUpload return new Upload object
func NewUpload() *Upload {
	return &Upload{}
}

// ConfigCommand configure kingpin options
func (u *Upload) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("upload", "Upload file to jutge.org").Action(u.Run)

	// Arguments
	cmd.Arg("file", "File to upload").Required().ExistingFilesVar(&u.files)

	// Flags
	cmd.Flag("compiler", "Compiler to use").Short('C').Default("G++11").StringVar(&u.compiler)
	cmd.Flag("code", "Problem code").Short('c').StringVar(&u.code)
}

// Run the command
func (u *Upload) Run(c *kingpin.ParseContext) error {
	var err error

	a, err := auth.Login()
	if err != nil {
		return err
	}

	extractCode := u.code == ""

	var wg sync.WaitGroup

	sem := make(chan bool, 3)

	for _, file := range u.files {
		sem <- true
		wg.Add(1)
		go func(f string) {
			defer func() { <-sem; wg.Done() }()
			fmt.Println("Uploading:", f)
			if extractCode {
				u.code, err = getCode(Conf.Regex, f)
				if err != nil {
					fmt.Println("Can't get code for file:", f)
					return
				}
			}
			_, err = u.uploadFile(f, a)
			if err != nil {
				fmt.Println("Upload failed", f, err)
			}
		}(file)
	}

	wg.Wait()

	return err
}

func (u Upload) uploadFile(fileName string, a auth.Auth) (*req.Resp, error) {

	param := req.Param{
		"annotation":  "hi",
		"compiler_id": u.compiler,
		"submit":      "submit",
		"token_uid":   a.TokenUID,
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file := req.FileUpload{
		File:      f,
		FieldName: "file",
		FileName:  "program",
	}

	return a.R.Post("https://jutge.org/problems/"+u.code+"/submissions", param, file)

}
