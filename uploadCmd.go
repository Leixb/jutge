package main

import (
	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Upload settings
type UploadCmd struct {
	files      []string
	code       string
	compiler   string
	annotation string
	check      bool
}

// ConfigCommand configure kingpin options
func (u *UploadCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("upload", "Upload file to jutge.org").Action(u.Run)

	// Arguments
	cmd.Arg("file", "File to upload").Required().ExistingFilesVar(&u.files)

	// Flags
	cmd.Flag("compiler", "Compiler to use").Short('C').Default("G++11").EnumVar(&u.compiler, commands.GetCompilers()...)
	cmd.Flag("code", "Problem code").Short('c').StringVar(&u.code)
	cmd.Flag("annotation", "Annotation").Short('a').Default("Uploaded with jutge_cli go").StringVar(&u.annotation)
	cmd.Flag("check", "Check veredict after upload").BoolVar(&u.check)
}

// Run the command
func (u *UploadCmd) Run(*kingpin.ParseContext) error {
	cmd := commands.NewUpload()
	cmd.Code = u.code
	cmd.Compiler = u.compiler
	cmd.Annotation = u.annotation
	cmd.Check = u.check

	err := cmd.UploadFiles(u.files)
	if err != nil {
		return err
	}
	if u.check {
		return cmd.CheckUploaded()
	}
	return nil
}
