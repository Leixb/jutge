package main

import (
	"fmt"

	"github.com/leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Test settings
type TestCmd struct {
	code     string
	programs []string
}

// ConfigCommand configure kingpin options
func (t *TestCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("test", "Test program").Action(t.Run)

	// Arguments
	cmd.Arg("programs", "Program to test").ExistingFilesVar(&t.programs)

	// Flags
	cmd.Flag("code", "Code of program to use").Short('c').StringVar(&t.code)
}

// Run the command
func (t *TestCmd) Run(c *kingpin.ParseContext) error {
	cmd := commands.NewTest()
	cmd.Code = t.code

	passed, count, err := cmd.TestPrograms(t.programs)
	if err != nil {
		return err
	}

	if len(t.programs) > 1 {
		fmt.Printf("=== Success: %d/%d\n", passed, count)
	}
	if passed != count {
		return fmt.Errorf("Failed %d out of %d tests", count-passed, count)
	}
	return nil
}
