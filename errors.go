package utils

import (
	"errors"
	"strings"
)

func HandleErrorMessage(err error) string {
	return strings.ReplaceAll(strings.Split(err.Error(), ";")[0], "mssql: ", "")
}

func HandleNewErrorMessage(err error) error {
	return errors.New(HandleErrorMessage(err))
}
