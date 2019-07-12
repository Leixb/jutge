package main

import (
	"os"

	"github.com/imroc/req"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge/auth"
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
	a         *auth.Credentials
}

var Conf GlobalConfig

func (c *GlobalConfig) getToken() (string, error) {
	if Conf.a == nil {
		Conf.a = auth.GetInstance()
	}
	return c.a.TokenUID, nil
}

func (c *GlobalConfig) getReq() (*req.Req, error) {
	if Conf.a == nil {
		Conf.a = auth.GetInstance()
	}
	return c.a.R, nil
}

func main() {
	app := kingpin.New(os.Args[0], "Jutge.org CLI implemented in go")

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
