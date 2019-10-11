package main

import (
	"fmt"

	"github.com/Leixb/jutge/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

type testCmd struct {
	code     string
	programs []string

	noDownload bool
}

func (t *testCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("test", "Test program").Action(t.Run)

	// Arguments
	cmd.Arg("programs", "Program to test").ExistingFilesVar(&t.programs)

	// Flags
	cmd.Flag("code", "Code of program to use").Short('c').StringVar(&t.code)
	cmd.Flag("no-download", "Don't attempt to download test files if not found in system").BoolVar(&t.noDownload)
}

// Run the command
func (t *testCmd) Run(c *kingpin.ParseContext) error {
	cmd := commands.NewTest()
	cmd.Code = t.code
	cmd.DownloadMissing = !t.noDownload

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
