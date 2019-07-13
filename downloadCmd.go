package main

import (
	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Download settings
type DownloadCmd struct {
	codes     []string
	overwrite bool
}

// ConfigCommand configure kingpin options
func (d *DownloadCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("download", "Download problem files from jutge.org").Action(d.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to download").Required().StringsVar(&d.codes)

	// Flags
	cmd.Flag("overwrite", "Overwrite existing files").BoolVar(&d.overwrite)
}

// Run the command
func (d *DownloadCmd) Run(*kingpin.ParseContext) error {
	cmd := commands.NewDownload()
	cmd.Overwrite = true
	return cmd.DownloadProblems(d.codes)
}
