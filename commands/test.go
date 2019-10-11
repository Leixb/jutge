package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type test struct {
	Code            string
	DownloadMissing bool
}

// NewTest returns test object
func NewTest() *test {
	return &test{Code: "", DownloadMissing: false}
}

// TestPrograms concurrently tests all the programs in `programs []string`
func (t *test) TestPrograms(programs []string) (passedTotal, countTotal int, err error) {
	for _, fileName := range programs {

		var code string

		if t.Code == "" {
			code, err = getCode(fileName)
			if err != nil {
				fmt.Println(" ! Can't get code for", fileName, err)
				continue
			}
		}

		passed, count, err := t.TestProgram(code, fileName)
		if err != nil {
			fmt.Println(" !  Error running tests for", fileName)
			continue
		}

		fmt.Printf(" #  %s Success: %d/%d\n", fileName, passed, count)

		passedTotal += passed
		countTotal += count

	}
	return
}

// TestProgram tests program fileName against all sample files for the given code found at Conf.WorkDir
// If there is no folder for the code, it tries to download the files from jutge.org (Downloading can be
// disabled by setting t.DownloadMissing to False).
func (t *test) TestProgram(code, fileName string) (passed, count int, err error) {
	folder := filepath.Join(conf.workDir, code)

	if _, err := os.Stat(folder); os.IsNotExist(err) && t.DownloadMissing {
		fmt.Println(" -", folder, "does not exist, downloading...")
		err = NewDownload().DownloadProblem(code)
		if err != nil {
			return 0, 0, err
		}
	}

	inputFiles, err := filepath.Glob(folder + "/*.inp")
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan bool, conf.concurrency)

	for _, inputFile := range inputFiles {

		count++

		sem <- true
		wg.Add(1)
		go func(iFile string) {
			defer func() { <-sem; wg.Done() }()

			ok, err := t.runTest(fileName, iFile)
			if err != nil {
				fmt.Println(" ! Error on", iFile, err)
			}

			if ok {
				passed++
			}

		}(inputFile)

	}

	wg.Wait()

	return
}

// runTest tests program against a single sample. If the output of the program
// does not match the expected output it prints an error and a diff
func (t *test) runTest(command, iFile string) (bool, error) {
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
		fmt.Println(" #  OK:", iFile)
		return true, nil
	}

	// Results don't match -> output diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(expected), string(output), true)

	str := fmt.Sprintf(" !  FAILED: %s\n", iFile)
	str = fmt.Sprintf("%s ===== OUTPUT =====\n%s\n", str, string(output))
	str = fmt.Sprintf("%s ==== EXPECTED ====\n%s\n", str, string(expected))
	str = fmt.Sprintf("%s ====== DIFF ======\n%s\n", str, dmp.DiffPrettyText(diffs))
	str = fmt.Sprintf("%s ==================\n", str)

	fmt.Print(str)
	return false, nil
}

// runCommand run command with input from file inputFile and return output
func (t *test) runCommand(command, inputFile string) ([]byte, error) {

	if len(command) == 0 {
		return nil, fmt.Errorf("Empty command")
	}

	if command[0] != '/' {
		command = "./" + command
	}

	cmd := exec.Command(command)

	input, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	cmd.Stdin = input

	output, err := cmd.CombinedOutput()

	return output, err
}
