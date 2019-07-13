package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge/commands"
)

type JutgeCommand interface {
	ConfigCommand(*kingpin.Application)
	Run(*kingpin.ParseContext) error
}

func main() {
	app := kingpin.New(os.Args[0], "Jutge.org CLI")

	kingpin.Flag("work-dir",
		"Directory to save jutge files").
		Default("JutgeProblems").
		Envar("JUTGE_WD").
		StringVar(commands.Regex())
	kingpin.Flag("concurrency",
		"Maximum concurrent routines").
		Default("3").
		UintVar(commands.Concurrency())

	commands := []JutgeCommand{
		&DownloadCmd{},
		&TestCmd{},
		&UploadCmd{},
		&CheckCmd{},
	}

	for _, command := range commands {
		command.ConfigCommand(app)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
