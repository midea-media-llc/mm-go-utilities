package utils

import (
	"errors"

	mssql "github.com/denisenkom/go-mssqldb"
	"google.golang.org/grpc/status"
)

func HandleSqlError(err error) error {
	e := err.(mssql.Error)
	return errors.New(e.SQLErrorMessage())
}

func HandleGrpcError(err error) error {
	if e, ok := status.FromError(err); ok {
		return errors.New(e.Message())
	}
	return err
}
