package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge/commands"
)

type jutgeCommand interface {
	ConfigCommand(*kingpin.Application)
	Run(*kingpin.ParseContext) error
}

var (
	username,
	password *string
)

func setUsername(*kingpin.ParseContext) error {
	commands.SetUsername(*username)
	return nil
}
func setPass(*kingpin.ParseContext) error {
	commands.SetPassword(*password)
	return nil
}

func main() {
	app := kingpin.New("jutge", "A jutge.org client in your terminal!").
		DefaultEnvars().
		Author("Leixb").
		Version("v0.3.0")

	app.Flag("work-dir",
		"Directory to save jutge files").
		Default("JutgeProblems").
		StringVar(commands.WorkDir())
	app.Flag("concurrency",
		"Maximum concurrent routines").
		Default("3").
		UintVar(commands.Concurrency())
	app.Flag("regex", "Regular expression used to validate and find problem codes in filenames").RegexpVar(commands.Regex())

	username = app.Flag("user", "Username").String()
	password = app.Flag("pass", "Password").String()

	app.Action(setUsername).Action(setPass)

	commands := []jutgeCommand{
		&downloadCmd{},
		&testCmd{},
		&uploadCmd{},
		&checkCmd{},
		&databaseCmd{},
		&newCmd{},
	}

	for _, command := range commands {
		command.ConfigCommand(app)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
