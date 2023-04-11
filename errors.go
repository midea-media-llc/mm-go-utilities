package utils

import (
	"errors"
	"net/http"
	"reflect"

	"google.golang.org/grpc/codes"
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
		if e.Code() == codes.InvalidArgument {
			return NewValidateHttpError(err)
		}
		return errors.New(e.Message())
	}
	return err
}

func NewValidateHttpError(err error) error {
	if e, ok := status.FromError(err); ok {
		if e.Code() == codes.InvalidArgument {
			return &validateError{err: "888", field: e.Message(), status: http.StatusBadRequest}
		}
		return errors.New(e.Message())
	}

	return err
}

type IValidateError interface {
	Error() string
	Status() int
	Field() string
}

type validateError struct {
	err    string
	field  string
	status int
}

func (s *validateError) Error() string {
	return s.err
}

func (s *validateError) Field() string {
	return s.field
}

func (s *validateError) Status() int {
	return s.status
}
