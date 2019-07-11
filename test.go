package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

// Test object that wraps its settings
type Test struct {
	code string
	file *os.File
}

// NewTest return new Test object
func NewTest() *Test {
	return &Test{}
}

// ConfigCommand configure kingpin options
func (t *Test) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("test", "Test program").Action(t.Run)

	// Arguments
	cmd.Arg("file", "Program to test").FileVar(&t.file)

	// Flags
	cmd.Flag("code", "Code of program to use").Short('c').StringVar(&t.code)
}

// Run the command
func (u *Test) Run(c *kingpin.ParseContext) error {
	return nil
}
