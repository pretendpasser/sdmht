package kitx

import (
	"io"
	"sync"
	"time"

	"sdmht/lib"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Response struct {
	Value interface{}
	Error error
}

func ServerEndpoint(makeEndpoint func() (endpoint.Endpoint, string), options *ServerOptions) endpoint.Endpoint {
	ep, name := makeEndpoint()

	if options.rateLimitOption.limiter != nil {
		ep = ratelimit.NewErroringLimiter(options.rateLimitOption.limiter)(ep)
	}
	if options.circuitBreakerOption.enable {
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	}
	if options.openTracingOption.otTracer != nil {
		ep = opentracing.TraceServer(options.openTracingOption.otTracer, name)(ep)
	}
	if options.zipkinTracerOption.tracer != nil {
		ep = zipkin.TraceEndpoint(options.zipkinTracerOption.tracer, name)(ep)
	}
	if options.Logger() != nil {
		ep = LoggingMiddleware(name, options.Logger())(ep)
	}
	if options.metricsOption.histogram != nil {
		ep = InstrumentingMiddleware(options.metricsOption.histogram.With("method", name))(ep)
	}

	return ErrorWrapperForServer(ep)
}

var GRPCConnections = make(map[string]*grpc.ClientConn)
var GRPCConnectionsMu sync.Mutex

func newGRPCClientFactory(makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts *ClientOptions) sd.Factory {
	return func(instance string) (i endpoint.Endpoint, closer io.Closer, e error) {

		GRPCConnectionsMu.Lock()
		conn, ok := GRPCConnections[instance]
		if !ok {
			var kacp = keepalive.ClientParameters{
				Time:                opts.keepaliveOption.time,
				Timeout:             opts.keepaliveOption.timeout,
				PermitWithoutStream: true, // send pings even without active streams
			}
			c, err := grpc.Dial(instance, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
			if err != nil {
				GRPCConnectionsMu.Unlock()
				return nil, nil, err
			}
			GRPCConnections[instance] = c
			conn = c
		}
		GRPCConnectionsMu.Unlock()

		ep, name := makeEndpoint(conn)

		if opts.openTracingOption.otTracer != nil {
			ep = opentracing.TraceClient(opts.openTracingOption.otTracer, name)(ep)
		}

		// if opts.rateLimitOption.enable {
		// 	ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(opts.rateLimitOption.every), opts.rateLimitOption.tokenCnt))(ep)
		// }

		// if opts.circuitBreakerOption.enable {
		// 	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		// 		Name:    name,
		// 		Timeout: opts.circuitBreakerOption.timeout,
		// 	}))(ep)
		// }

		return ep, nil, nil
	}
}

func GRPCClientEndpoint(instancer sd.Instancer, makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts *ClientOptions) endpoint.Endpoint {
	factory := newGRPCClientFactory(makeEndpoint, opts)

	logger := opts.Logger()
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	retryMax := 1
	timeout := 3 * time.Second
	if opts.loadBalanceOption.retryMax > 0 {
		retryMax = opts.loadBalanceOption.retryMax
	}
	if opts.loadBalanceOption.timeout > 0 {
		timeout = opts.loadBalanceOption.timeout
	}

	ep := lb.RetryWithCallback(timeout, balancer, func(n int, received error) (bool, error) {
		if _, ok := received.(lib.Error); ok {
			return false, received
		}
		return n < retryMax, received
	})
	return ErrorWrapperForClient(ep)
}
