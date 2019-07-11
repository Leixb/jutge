package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

// Download object that wraps its settings
type Download struct {
	code      string
	overwrite bool
}

// NewDownload return new Download object
func NewDownload() *Download {
	return &Download{}
}

// ConfigCommand configure kingpin options
func (d *Download) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("dowload", "Download problem files from jutge.org").Action(d.Run)

	// Arguments
	cmd.Arg("code", "Code of problem to download").Required().StringVar(&d.code)

	// Flags
	cmd.Flag("overwrite", "Overwrite existing files").BoolVar(&d.overwrite)
}

// Run the command
func (u *Download) Run(c *kingpin.ParseContext) error {
	return nil
}
