package server

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.Violations),
	}
}

func protovalidateVialationsToGoogleViolations(vs []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldPath,
			Description: v.Message,
		}
	}
	return res
}
