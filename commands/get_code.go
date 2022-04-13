package commands

import (
	"errors"
)

// getCode match regex against string and return matching code
func (j *jutge) getCode(fileName string) (string, error) {
	code := j.regex.FindString(fileName)
	if len(code) == 0 {
		return "", errors.New("No match")
	}

	return code, nil
}

// getCodeOrSame getCode or return the original value if not matched
func (j *jutge) getCodeOrSame(fileName string) string {
	code, err := j.getCode(fileName)
	if err != nil {
		return fileName
	}
	return code
}
