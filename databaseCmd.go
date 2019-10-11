package main

import (
	"fmt"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Leixb/jutge/commands"
	"github.com/Leixb/jutge/database"
)

type databaseCmd struct {
	code, title, zipFile string
	DBFile               string
}

func (d *databaseCmd) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("db", "Manage local database (use with caution)").Action(d.Run)

	//SubCommands
	cmd.Command("print", "Print contents of database").Action(d.printRun)

	addCmd := cmd.Command("add", "Add entry to database").Action(d.addRun)
	addCmd.Arg("code", "Proble ID").Required().StringVar(&d.code)
	addCmd.Arg("title", "Proble title").Required().StringVar(&d.title)

	queryCmd := cmd.Command("query", "Query title from database").Action(d.queryRun)
	queryCmd.Arg("code", "Proble ID").Required().StringVar(&d.code)

	importCmd := cmd.Command("import", "Import data from zip (this is quite usless atm)").Action(d.importRun)
	importCmd.Arg("zipFile", "Zip file").Required().StringVar(&d.zipFile)

}

func (d *databaseCmd) Run(*kingpin.ParseContext) error {
	d.DBFile = filepath.Join(*commands.WorkDir(), "jutge.db")
	return nil
}

func (d *databaseCmd) printRun(*kingpin.ParseContext) error {
	database.NewJutgeDB(d.DBFile).Print()
	return nil
}
func (d *databaseCmd) addRun(*kingpin.ParseContext) error {
	return database.NewJutgeDB(d.DBFile).Add(d.code, d.title)
}
func (d *databaseCmd) queryRun(*kingpin.ParseContext) error {
	title, err := database.NewJutgeDB(d.DBFile).Query(d.code)
	fmt.Println(title)
	return err
}
func (d *databaseCmd) importRun(*kingpin.ParseContext) error {

	return database.NewJutgeDB(d.DBFile).ImportZip(d.zipFile)
}
