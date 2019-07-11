package main

import (
	"os"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge_go/auth"
)

// Upload object that wraps its settings
type Upload struct {
	compiler string
	file     *os.File
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
	cmd.Arg("file", "File to upload").Required().FileVar(&u.file)

	// Flags
	cmd.Flag("compiler", "Compiler to use").Default("G++11").StringVar(&u.compiler)
	cmd.Flag("code", "Problem code").StringVar(&u.code)
}

// Run the command
func (u *Upload) Run(c *kingpin.ParseContext) error {
	var err error
	if u.code == "" {
		u.code, err = getCode(Conf.Regex, u.file.Name())
		if err != nil {
			return err
		}
	}

	a, err := auth.Login()
	if err != nil {
		return err
	}

	_, err = u.uploadFile(a)

	return err
}

func (u Upload) uploadFile(a auth.Auth) (*req.Resp, error) {

	param := req.Param{
		"annotation":  "hi",
		"compiler_id": u.compiler,
		"submit":      "submit",
		"token_uid":   a.TokenUID,
	}

	file := req.FileUpload{
		File:      u.file,
		FieldName: "file",
		FileName:  u.file.Name(),
	}

	return a.R.Post("https://jutge.org/problems/"+u.code+"/submissions", param, file)

}
