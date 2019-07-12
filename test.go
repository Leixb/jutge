package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Test object that wraps its settings
type Test struct {
	code string
	file *os.File

	concurrency int
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
	cmd.Flag("concurrency", "Number of simultaneous tests").Default("3").IntVar(&t.concurrency)
}

// Run the command
func (t *Test) Run(c *kingpin.ParseContext) error {

	var err error
	if t.code == "" {
		t.code, err = getCode(Conf.Regex, t.file.Name())
		if err != nil {
			return err
		}
	}

	folder := filepath.Join(Conf.WorkDir, t.code)

	files, err := filepath.Glob(folder + "/*.inp")
	if err != nil {
		return err
	}

	sem := make(chan bool, t.concurrency)

	var wg sync.WaitGroup

	for _, inputFile := range files {

		sem <- true
		wg.Add(1)
		go func(iFile string) {
			err = t.runTest(iFile)
			if err != nil {
				fmt.Println("Error on", iFile, err)
			}
			wg.Done()
			<-sem
		}(inputFile)

	}
	wg.Wait()
	return nil
}

func (t *Test) runTest(iFile string) error {
	output, err := t.runCommand(iFile)
	if err != nil {
		return err
	}

	expectedOutputFile := strings.TrimSuffix(iFile, filepath.Ext(iFile)) + ".cor"

	expected, err := ioutil.ReadFile(expectedOutputFile)
	if err != nil {
		return err
	}

	if bytes.Equal(output, expected) {
		fmt.Println("=== OK:", iFile)
	} else {

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expected), string(output), true)

		str := fmt.Sprintf("=== FAILED: %s\n", iFile)
		str = fmt.Sprintf("%s===== OUTPUT =====\n%s\n", str, string(output))
		str = fmt.Sprintf("%s==== EXPECTED ====\n%s\n", str, string(expected))
		str = fmt.Sprintf("%s====== DIFF ======\n%s\n", str, dmp.DiffPrettyText(diffs))
		str = fmt.Sprintf("%s==================\n", str)

		fmt.Print(str)

	}
	return nil
}

func (t *Test) runCommand(inputFile string) ([]byte, error) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("./" + t.file.Name())

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(stdin, string(input))
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return output, nil
}
