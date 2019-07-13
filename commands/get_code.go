package commands

import (
	"errors"
	"regexp"
)

// getCode match regex agains string and return matching code
func getCode(fileName string) (string, error) {
	re, err := regexp.Compile(conf.regex)
	if err != nil {
		return "", err
	}

	code := re.FindString(fileName)
	if len(code) == 0 {
		return "", errors.New("No match")
	}

	return code, nil
}

// getCodeOrSame getCode or return the original value if not matched
func getCodeOrSame(fileName string) string {
	code, err := getCode(fileName)
	if err != nil {
		return fileName
	}
	return code
}
