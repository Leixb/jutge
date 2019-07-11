package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Upload object that wraps its settings
type Upload struct {
	compiler string
	file     *os.File
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
}

// Run the command
func (u *Upload) Run(c *kingpin.ParseContext) error {
	return nil
}
