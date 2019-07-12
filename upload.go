package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Upload object that wraps its settings
type Upload struct {
	compiler    string
	files       []string
	code        string
	annotation  string
	concurrency int
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
	cmd.Flag("annotation", "Annotation").Short('a').Default("Uploaded with jutge_cli go").StringVar(&u.annotation)
	cmd.Flag("concurrency", "Number of simultaneous uploads").Default("3").IntVar(&u.concurrency)
}

// Run the command
func (u *Upload) Run(*kingpin.ParseContext) error {
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
			err = u.uploadFile(f)
			if err != nil {
				fmt.Println("Upload failed", f, err)
			}
		}(file)
	}

	wg.Wait()

	return err
}

func (u Upload) uploadFile(fileName string) error {

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
