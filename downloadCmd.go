package main

import (
	"github.com/Leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

type downloadCmd struct {
	codes     []string
	overwrite bool
}

func (d *downloadCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("download", "Download problem files from jutge").
		Alias("down").
		Action(d.Run)

	// Arguments
	cmd.Arg("code", "Code(s) of problem(s) to download").Required().StringsVar(&d.codes)

	// Flags
	cmd.Flag("overwrite", "Overwrite existing files").BoolVar(&d.overwrite)
}

func (d *downloadCmd) Run(*kingpin.ParseContext) error {
	cmd := commands.NewDownload()
	cmd.Overwrite = true
	return cmd.DownloadProblems(d.codes)
}
