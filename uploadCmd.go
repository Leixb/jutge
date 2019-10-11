package main

import (
	"fmt"

	"github.com/Leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

type uploadCmd struct {
	files []string

	code       string
	compiler   string
	annotation string
	check      bool
}

func (u *uploadCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("upload", "Upload file to jutge.org").
		Alias("up").Alias("submit").
		Action(u.Run)

	// Arguments
	cmd.Arg("file", "File to upload").Required().ExistingFilesVar(&u.files)

	// Flags
	cmd.Flag("compiler", "Compiler to use").Short('C').Default("G++11").EnumVar(&u.compiler, commands.GetCompilers()...)
	cmd.Flag("code", "Problem code").Short('c').StringVar(&u.code)
	cmd.Flag("annotation", "Annotation").Short('a').Default("Uploaded with jutge_cli go").StringVar(&u.annotation)
	cmd.Flag("check", "Check veredict after upload").BoolVar(&u.check)
}

func (u *uploadCmd) Run(*kingpin.ParseContext) error {
	cmd := commands.NewUpload()
	cmd.Code = u.code
	cmd.Compiler = u.compiler
	cmd.Annotation = u.annotation
	cmd.Check = u.check

	err := cmd.UploadFiles(u.files)
	if err != nil {
		return err
	}
	fmt.Println("=== Upload Finished")
	if u.check {
		fmt.Println("=== Waiting for veredicts...")
		return cmd.CheckUploaded()
	}
	return nil
}
