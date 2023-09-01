package server

import (
	"context"
	"sdmht/account/api"
	"sdmht/account/api/grpc/pb"
	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib/kitx"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
)

var _ pb.AccountServer = (*grpcServer)(nil)

type grpcServer struct {
	*pb.UnimplementedAccountServer

	registerHandler     grpc.Handler
	loginHandler        grpc.Handler
	logoutHandler       grpc.Handler
	authenticateHandler grpc.Handler
	getAccountHandler   grpc.Handler
}

func NewGRPCServer(svc itfs.Service, opts *kitx.ServerOptions) pb.AccountServer {
	srv := &grpcServer{}

	logger := opts.Logger()

	options := []grpc.ServerOption{
		grpc.ServerBefore(opts.MetadataToCtx("sdmht-account")),
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	srv.authenticateHandler = makeAuthenticateHandler(svc, options, opts)
	srv.registerHandler = makeRegisterHandler(svc, options, opts)
	srv.loginHandler = makeLoginHandler(svc, options, opts)
	srv.logoutHandler = makeLogoutHandler(svc, options, opts)
	srv.getAccountHandler = makeGetAccountHandler(svc, options, opts)

	return srv
}

func (s *grpcServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	_, rsp, err := s.registerHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.(*pb.RegisterRes), nil
}

func (s *grpcServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	_, rsp, err := s.loginHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.(*pb.LoginRes), nil
}

func (s *grpcServer) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutRes, error) {
	_, rsp, err := s.logoutHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.(*pb.LogoutRes), nil
}

func (s *grpcServer) Authenticate(ctx context.Context, req *pb.AuthenticateReq) (*pb.AuthenticateRes, error) {
	_, rsp, err := s.authenticateHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.(*pb.AuthenticateRes), nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountReq) (*pb.GetAccountRes, error) {
	_, rsp, err := s.getAccountHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp.(*pb.GetAccountRes), nil
}

func makeRegisterHandler(svc itfs.Service, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeRegisterEndpoint(svc)
		return ep, "account_svc.Register"
	}, opts)
	return grpc.NewServer(ep, decodeRegisterReq, enRegisterRes, options...)
}

func makeLoginHandler(svc itfs.Service, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeLoginEndpoint(svc)
		return ep, "account_svc.Login"
	}, opts)
	return grpc.NewServer(ep, decodeLoginReq, enLoginRes, options...)
}

func makeLogoutHandler(svc itfs.Service, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeLogoutEndpoint(svc)
		return ep, "account_svc.Logout"
	}, opts)
	return grpc.NewServer(ep, decodeLogoutReq, enLogoutRes, options...)
}

func makeAuthenticateHandler(svc itfs.Service, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeAuthenticateEndpoint(svc)
		return ep, "account_svc.Authenticate"
	}, opts)
	return grpc.NewServer(ep, decodeAuthenticateReq, enAuthenticateRes, options...)
}

func makeGetAccountHandler(svc itfs.Service, options []grpc.ServerOption, opts *kitx.ServerOptions) grpc.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeGetAccountEndpoint(svc)
		return ep, "account_svc.GetAccount"
	}, opts)
	return grpc.NewServer(ep, decodeGetAccountReq, enGetAccountRes, options...)
}
