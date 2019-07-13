package main

import (
	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Check settings
type CheckCmd struct {
	codes []string
}

// ConfigCommand configure kingpin options
func (c *CheckCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("check", "Check problem files from jutge.org").Action(c.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to check").Required().StringsVar(&c.codes)

}

// Run the command
func (c *CheckCmd) Run(*kingpin.ParseContext) error {
	return commands.NewCheck().CheckProblems(c.codes)
}
