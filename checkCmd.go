package main

import (
	"fmt"

	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

type checkCmd struct {
	codes []string

	submission int
}

func (c *checkCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("check", "Check problem files from jutge.org").Action(c.Run)

	// Arguments
	cmd.Arg("code", "Codes of problems to check").Required().StringsVar(&c.codes)

	// Flags
	cmd.Flag("submission", "Check submission by number (if 0 checks last submission)").Default("-1").PlaceHolder("2").Short('s').IntVar(&c.submission)

}

func (c *checkCmd) Run(*kingpin.ParseContext) error {
	if c.submission != -1 {
		if len(c.codes) != 1 {
			return fmt.Errorf("submission only accepts 1 code to check")
		}
		if c.submission == 0 {
			veredict, err := commands.NewCheck().CheckLast(c.codes[0])
			fmt.Printf(" - %s last: %s\n", c.codes[0], veredict)
			return err
		}
		veredict, err := commands.NewCheck().CheckSubmission(c.codes[0], c.submission)
		fmt.Printf(" - %s S%03d: %s\n", c.codes[0], c.submission, veredict)
		return err
	}
	return commands.NewCheck().CheckProblems(c.codes)
}
