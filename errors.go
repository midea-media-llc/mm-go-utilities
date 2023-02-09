package utils

import (
	"errors"
	"strings"

	"google.golang.org/grpc/status"
)

func HandleErrorMessage(err error) string {
	return strings.ReplaceAll(strings.Split(err.Error(), ";")[0], "mssql: ", "")
}

func HandleNewErrorMessage(err error) error {
	return errors.New(HandleErrorMessage(err))
}

func HandleGrpcError(err error) error {
	if e, ok := status.FromError(err); ok {
		return errors.New(e.Message())
	}
	return err
}
