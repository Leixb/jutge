package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Test settings
type Test struct {
	code     string
	programs []string

	concurrency int
}

// NewTest return Test object
func NewTest() *Test {
	return &Test{concurrency: 3}
}

// ConfigCommand configure kingpin options
func (t *Test) ConfigCommand(app *kingpin.Application) {
	cmd := app.Command("test", "Test program").Action(t.Run)

	// Arguments
	cmd.Arg("programs", "Program to test").ExistingFilesVar(&t.programs)

	// Flags
	cmd.Flag("code", "Code of program to use").Short('c').StringVar(&t.code)
	cmd.Flag("concurrency", "Number of simultaneous tests").Default("3").IntVar(&t.concurrency)
}

// Run the command
func (t *Test) Run(c *kingpin.ParseContext) error {
	passed, count, err := t.TestPrograms()
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

// TestPrograms Test all the programs in t.programs
func (t *Test) TestPrograms() (passedTotal, countTotal int, err error) {
	for _, fileName := range t.programs {

		var code string

		if t.code == "" {
			code, err = getCode(fileName)
			if err != nil {
				fmt.Println("Can't get error for", fileName, err)
				continue
			}
		}

		passed, count, err := t.TestProgram(code, fileName)
		if err != nil {
			fmt.Println("=== Error running tests for", fileName)
			continue
		}

		fmt.Printf("=== %s Success: %d/%d\n", fileName, passed, count)

		passedTotal += passed
		countTotal += count

	}
	return
}

// TestProgram Test program fileName against all sample files in Conf.WorkDir
func (t *Test) TestProgram(code, fileName string) (passed, count int, err error) {
	folder := filepath.Join(Conf.WorkDir, t.code)

	inputFiles, err := filepath.Glob(folder + "/*.inp")
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan bool, t.concurrency)

	for _, inputFile := range inputFiles {

		count++

		sem <- true
		wg.Add(1)
		go func(iFile string) {
			defer func() { <-sem; wg.Done() }()

			ok, err := t.runTest(fileName, iFile)
			if err != nil {
				fmt.Println("Error on", iFile, err)
			}

			if ok {
				passed++
			}

		}(inputFile)

	}

	wg.Wait()

	return
}

// runTest test program against a single sample
func (t *Test) runTest(command, iFile string) (bool, error) {
	output, err := t.runCommand(command, iFile)
	if err != nil {
		return false, err
	}

	expectedOutputFile := strings.TrimSuffix(iFile, filepath.Ext(iFile)) + ".cor"

	expected, err := ioutil.ReadFile(expectedOutputFile)
	if err != nil {
		return false, err
	}

	if bytes.Equal(output, expected) {
		fmt.Println("=== OK:", iFile)
		return true, nil
	} else {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expected), string(output), true)

		str := fmt.Sprintf("=== FAILED: %s\n", iFile)
		str = fmt.Sprintf("%s===== OUTPUT =====\n%s\n", str, string(output))
		str = fmt.Sprintf("%s==== EXPECTED ====\n%s\n", str, string(expected))
		str = fmt.Sprintf("%s====== DIFF ======\n%s\n", str, dmp.DiffPrettyText(diffs))
		str = fmt.Sprintf("%s==================\n", str)

		fmt.Print(str)
		return false, nil
	}
}

// runCommand run command with input from file inputFile and return output
func (t *Test) runCommand(command, inputFile string) ([]byte, error) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("./" + command)

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
