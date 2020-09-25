package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Leixb/jutge/commands"
)

type newCmd struct {
	code, ext string
	dryRun    bool
}

func (n *newCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("new", "Create a new file for a problem in the current directory").
		Action(n.Run)

	// Arguments
	cmd.Arg("code", "Code of problem").Required().StringVar(&n.code)
	cmd.Arg("ext", "Extension of file").Default("cpp").StringVar(&n.ext)

	// Flags
	cmd.Flag("dry-run", "Only print filename, do not create file").BoolVar(&n.dryRun)
}

func (n *newCmd) Run(*kingpin.ParseContext) error {
	cmd := commands.NewNewfile()
	cmd.Code = n.code
	cmd.Extension = n.ext

	filename, err := cmd.GetFilename()
	if err != nil {
		return err
	}

	if !n.dryRun {
		os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	}
	fmt.Println(filename)

	return nil
}
