package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// TestPrograms concurrently tests all the programs in `programs []string`
func (j *jutge) TestPrograms(globalCode string, programs []string, downloadMissing, overwrite bool, concurrency uint) (passedTotal, countTotal int, err error) {
	for _, fileName := range programs {

		code := globalCode

		if globalCode == "" {
			code, err = getCode(fileName, j.regex)
			if err != nil {
				fmt.Println(" ! Can't get code for", fileName, err)
				continue
			}
		}

		passed, count, err := j.TestProgram(code, fileName, downloadMissing, overwrite, concurrency)
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
func (j *jutge) TestProgram(code, fileName string, downloadMissing bool, overwrite bool, concurrency uint) (passed, count int, err error) {
	folder := filepath.Join(j.folder, code)

	if _, err := os.Stat(folder); os.IsNotExist(err) && downloadMissing {
		fmt.Println(" -", folder, "does not exist, downloading...")
		err = j.DownloadProblem(code, j.folder, overwrite)
		if err != nil {
			return 0, 0, err
		}
	}

	inputFiles, err := filepath.Glob(folder + "/*.inp")
	if err != nil {
		return
	}

	passed = 0
	failed := 0
	resultChan := make(chan bool)

	go func() {
		for ok := range resultChan {
			if ok {
				passed++
			} else {
				failed++
			}
		}
	}()

	RunParallelFuncs(inputFiles, func(iFile string) error {
		ok, err := runTest(fileName, iFile)
		if err != nil {
			return err
		}
		resultChan <- ok
		return nil
	}, concurrency)

	close(resultChan)

	return passed, len(inputFiles), nil
}

// runTest tests program against a single sample. If the output of the program
// does not match the expected output it prints an error and a diff
func runTest(command, iFile string) (bool, error) {
	output, err := runCommand(command, iFile)
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
	diffs := dmp.DiffMain(string(output), string(expected), true)

	str := fmt.Sprintf(" !  FAILED: %s\n", iFile)
	str = fmt.Sprintf("%s ===== OUTPUT =====\n%s\n", str, string(output))
	str = fmt.Sprintf("%s ==== EXPECTED ====\n%s\n", str, string(expected))
	str = fmt.Sprintf("%s ====== DIFF ======\n%s\n", str, dmp.DiffPrettyText(diffs))
	str = fmt.Sprintf("%s ==================\n", str)

	fmt.Print(str)
	return false, nil
}

// runCommand run command with input from file inputFile and return output
func runCommand(command, inputFile string) ([]byte, error) {

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
