package grpcresponse

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error(err error, code codes.Code, message string) error {
	return status.Errorf(
		code,
		errors.WithMessage(err, message).Error(),
	)
}
