package kitx

import (
	"context"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc/metadata"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type contextKey string
type MetadataKey string

const SourceKey contextKey = "source"

var (
	DefaultKeepaliveOption = keepaliveOption{time: 5 * time.Second, timeout: 1 * time.Second}
)

type Option interface {
	apply(*options)
}

type options struct {
	circuitBreakerOption
	rateLimitOption
	openTracingOption
	zipkinTracerOption
	loggerOption
	loadBalanceOption
	metricsOption
	sourceOption
	metadataOption
	keepaliveOption
}

type circuitBreakerOption struct {
	enable  bool
	timeout time.Duration
}

func (o circuitBreakerOption) apply(opts *options) {
	opts.circuitBreakerOption = o
}

func WithCircuitBreaker(timeout time.Duration) Option {
	return circuitBreakerOption{false, timeout} // FIXME
}

type rateLimitOption struct {
	limiter ratelimit.Allower
}

func (o rateLimitOption) apply(opts *options) {
	opts.rateLimitOption = o
}

func WithRateLimit(limiter ratelimit.Allower) Option {
	return rateLimitOption{limiter}
}

type openTracingOption struct {
	otTracer stdopentracing.Tracer
}

func (o openTracingOption) apply(opts *options) {
	opts.openTracingOption = o
}

func WithOpenTracing(otTracer stdopentracing.Tracer) Option {
	return openTracingOption{otTracer}
}

type loggerOption struct {
	logger log.Logger
}

func (o loggerOption) apply(opts *options) {
	opts.loggerOption = o
}

func WithLogger(logger log.Logger) Option {
	return loggerOption{logger}
}

type loadBalanceOption struct {
	retryMax int
	timeout  time.Duration
}

func (o loadBalanceOption) apply(opts *options) {
	opts.loadBalanceOption = o
}

func WithLoadBalance(retryMax int, timeout time.Duration) Option {
	return loadBalanceOption{retryMax, timeout}
}

type metricsOption struct {
	histogram metrics.Histogram
}

func (o metricsOption) apply(opts *options) {
	opts.metricsOption = o
}

func WithMetrics(histogram metrics.Histogram) Option {
	return metricsOption{histogram}
}

type zipkinTracerOption struct {
	tracer *zipkin.Tracer
}

func (o zipkinTracerOption) apply(opts *options) {
	opts.zipkinTracerOption = o
}

func WithZipkinTracer(tracer *zipkin.Tracer) Option {
	return zipkinTracerOption{tracer}
}

type sourceOption struct {
	source string
}

func (o sourceOption) apply(opts *options) {
	opts.sourceOption = o
}

func WithSource(source string) Option {
	return sourceOption{source}
}

type metadataOption struct {
	md map[string][]string
}

func (o metadataOption) apply(opts *options) {
	opts.metadataOption = o
}

func WithMetadata(md map[string][]string) Option {
	return metadataOption{md}
}

type keepaliveOption struct {
	time    time.Duration // send pings every <time> if there is no activity
	timeout time.Duration // wait <timeout> for ping ack before considering the connection dead
}

func (o keepaliveOption) apply(opts *options) {
	opts.keepaliveOption = o
}

func (o keepaliveOption) Time() time.Duration {
	return o.time
}

func (o keepaliveOption) Timeout() time.Duration {
	return o.timeout
}

func WithKeepaliveOption(time, timeout time.Duration) Option {
	return keepaliveOption{time, timeout}
}

type ServerOptions struct {
	options
}

func NewServerOptions(opts ...Option) *ServerOptions {
	so := &ServerOptions{}
	for _, o := range opts {
		o.apply(&so.options)
	}
	return so
}

func (o *ServerOptions) Logger() log.Logger {
	return o.loggerOption.logger
}

func (o *ServerOptions) ZipkinTracer() *zipkin.Tracer {
	return o.zipkinTracerOption.tracer
}

func (o *ServerOptions) SourceToCtx() grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		// capital "Key" is illegal in HTTP/2.
		sourceHeader, ok := md[string(SourceKey)]
		if !ok || len(sourceHeader) != 1 {
			return ctx
		}

		ctx = context.WithValue(ctx, SourceKey, sourceHeader[0])

		return ctx
	}
}

// use 'namespace' to distinguish invoker and grpc's own metadata.
func (o *ServerOptions) MetadataToCtx(namespace string) grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		// capital "Key" is illegal in HTTP/2.
		for k, v := range md {
			if strings.HasPrefix(k, namespace) {
				ctx = context.WithValue(ctx, MetadataKey(k[len(namespace)+1:]), v)
			}
		}
		return ctx
	}
}

type ClientOptions struct {
	options
}

func NewClientOptions(opts ...Option) *ClientOptions {
	co := &ClientOptions{
		options: options{
			keepaliveOption: DefaultKeepaliveOption,
		},
	}
	for _, o := range opts {
		o.apply(&co.options)
	}
	return co
}

func (o *ClientOptions) Logger() log.Logger {
	return o.loggerOption.logger
}

func (o *ClientOptions) ZipkinTracer() *zipkin.Tracer {
	return o.zipkinTracerOption.tracer
}

func (o *ClientOptions) SourceToGRPC() grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		(*md)[string(SourceKey)] = []string{o.sourceOption.source}
		return ctx
	}
}

// use 'namespace' to distinguish invoker and grpc's own metadata.
func (o *ClientOptions) MetadataToGRPC(namespace string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		for k, v := range o.metadataOption.md {
			(*md)[namespace+"_"+k] = v
		}
		return ctx
	}
}
