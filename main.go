package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type JutgeCommand interface {
	ConfigCommand(*kingpin.Application)
	Run(*kingpin.ParseContext) error
}

type GlobalConfig struct {
	WorkDir   string
	Regex     string
	Verbosity int
	Quiet     bool
}

var Conf GlobalConfig

func main() {
	app := kingpin.New("jutge_go", "Jutge.org CLI implemented in go")

	app.Flag("work-dir", "Directory to save jutge files").Default("JutgeProblems").Envar("JUTGE_WD").StringVar(&Conf.WorkDir)
	app.Flag("regex", "Code regex").Default(`[PGQX]\d{5}_(ca|en|es)`).StringVar(&Conf.Regex)
	app.Flag("verbosity", "Verbosity level").Short('v').CounterVar(&Conf.Verbosity)
	app.Flag("quiet", "Suppress output").Short('q').BoolVar(&Conf.Quiet)

	commands := []JutgeCommand{
		NewDownload(),
		NewTest(),
		NewUpload(),
	}

	for _, command := range commands {
		command.ConfigCommand(app)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
