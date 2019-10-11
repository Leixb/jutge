package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge/commands"
	"github.com/leixb/jutge/database"

	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func normalizeString(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	res, _, _ := transform.String(t, s)
	res = strings.ReplaceAll(res, " ", "_")
	return res
}

type newCmd struct {
	code, ext string
	dryRun    bool
}

func (n *newCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("new", "Create file for problem").
		Action(n.Run)

	// Arguments
	cmd.Arg("code", "Code of problem").Required().StringVar(&n.code)
	cmd.Arg("ext", "Extension of file").Default("cpp").StringVar(&n.ext)

	// Flags
	cmd.Flag("dry-run", "Only print filename, do not create file").BoolVar(&n.dryRun)
}

func (n *newCmd) Run(*kingpin.ParseContext) error {
	dbFile := filepath.Join(*commands.WorkDir(), "jutge.db")
	problemName, err := database.NewJutgeDB(dbFile).Query(n.code)
	if err != nil {
		return err
	}

	problemName = normalizeString(problemName)

	filename := fmt.Sprintf("%s_%s.%s", n.code, problemName, n.ext)

	if !n.dryRun {
		os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	}
	fmt.Println(filename)

	return nil
}
