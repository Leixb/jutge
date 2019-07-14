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
	app := kingpin.New("jutge", "Jutge.org CLI").
		DefaultEnvars().
		Author("Leixb").
		Version("v0.2.0")

	kingpin.Flag("work-dir",
		"Directory to save jutge files").
		Default("JutgeProblems").
		Envar("JUTGE_WD").
		StringVar(commands.WorkDir())
	kingpin.Flag("concurrency",
		"Maximum concurrent routines").
		Default("3").
		UintVar(commands.Concurrency())

	commands := []JutgeCommand{
		&downloadCmd{},
		&testCmd{},
		&uploadCmd{},
		&checkCmd{},
	}

	for _, command := range commands {
		command.ConfigCommand(app)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
