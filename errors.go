package utils

import (
	"errors"
	"reflect"

	"google.golang.org/grpc/status"
)

type ISqlError interface {
	SQLErrorMessage() string
}

// HandleSqlError handles SQL errors by returning a new error object that contains the SQL error message.
// If the provided error is not of type mssql.Error, the function simply returns the original error.
func HandleSqlError(err error) error {
	if reflect.TypeOf(err) != TYPE_SQL_ERROR {
		return err
	}

	e := err.(ISqlError)
	return errors.New(e.SQLErrorMessage())
}

// HandleGrpcError handles gRPC errors by returning a new error object that contains the gRPC error message.
// If the provided error is not a gRPC error, the function simply returns the original error.
func HandleGrpcError(err error) error {
	if e, ok := status.FromError(err); ok {
		return errors.New(e.Message())
	}
	return err
}
