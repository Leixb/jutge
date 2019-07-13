package main

import (
	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

type checkCmd struct {
	codes []string
}

func (c *checkCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("check", "Check problem files from jutge.org").Action(c.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to check").Required().StringsVar(&c.codes)

}

func (c *checkCmd) Run(*kingpin.ParseContext) error {
	return commands.NewCheck().CheckProblems(c.codes)
}
