package commands

import (
	"errors"
	"regexp"
)

// getCode match regex against string and return matching code
func getCode(fileName string, regex *regexp.Regexp) (string, error) {
	code := regex.FindString(fileName)
	if len(code) == 0 {
		return "", errors.New("No match")
	}

	return code, nil
}

// getCodeOrSame getCode or return the original value if not matched
func getCodeOrSame(fileName string, regex *regexp.Regexp) string {
	code, err := getCode(fileName, regex)
	if err != nil {
		return fileName
	}
	return code
}
