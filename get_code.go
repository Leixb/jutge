package main

import (
	"errors"
	"regexp"
)

func getCode(regex, fileName string) (string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}

	code := re.FindString(fileName)
	if len(code) == 0 {
		return "", errors.New("No match")
	}

	return code, nil
}
