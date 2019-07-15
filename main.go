package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leixb/jutge/commands"
)

type JutgeCommand interface {
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
	app := kingpin.New("jutge", "Jutge.org CLI").
		DefaultEnvars().
		Author("Leixb").
		Version("v0.2.0")

	app.Flag("work-dir",
		"Directory to save jutge files").
		Default("JutgeProblems").
		StringVar(commands.WorkDir())
	app.Flag("concurrency",
		"Maximum concurrent routines").
		Default("3").
		UintVar(commands.Concurrency())
	app.Flag("regex", "Regular expression to match code").RegexpVar(commands.Regex())

	username = app.Flag("user", "Username").String()
	password = app.Flag("pass", "Password").String()

	app.Action(setUsername).Action(setPass)

	commands := []JutgeCommand{
		&downloadCmd{},
		&testCmd{},
		&uploadCmd{},
		&checkCmd{},
	}

	for _, command := range commands {
		command.ConfigCommand(app)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
