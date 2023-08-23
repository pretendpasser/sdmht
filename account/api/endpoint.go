package api

import (
	"context"

	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib/kitx"

	"github.com/go-kit/kit/endpoint"
)

func MakeAuthenticateEndpoint(s itfs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := request.(string)
		v, err := s.Authenticate(ctx, token)
		return kitx.Response{Value: v, Error: err}, nil
	}
}
