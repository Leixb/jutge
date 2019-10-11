package commands

import (
	"errors"
)

// getCode match regex against string and return matching code
func getCode(fileName string) (string, error) {
	code := conf.regex.FindString(fileName)
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
