package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	unexpectedError error = status.Errorf(codes.Internal, "unexpected error")
)
