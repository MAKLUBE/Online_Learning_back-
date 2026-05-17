package grpcerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrUnauthenticated = errors.New("unauthenticated")
)

func InvalidArgument(message string) error {
	return status.Error(codes.InvalidArgument, message)
}

func NotFound(message string) error {
	return status.Error(codes.NotFound, message)
}

func AlreadyExists(message string) error {
	return status.Error(codes.AlreadyExists, message)
}

func Unauthenticated(message string) error {
	return status.Error(codes.Unauthenticated, message)
}

func PermissionDenied(message string) error {
	return status.Error(codes.PermissionDenied, message)
}

func Internal(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	return status.Error(codes.Internal, err.Error())
}
