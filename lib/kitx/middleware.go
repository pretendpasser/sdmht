package kitx

import (
	"context"
	"fmt"
	"net/http"
	"sdmht/lib"
	"sdmht/lib/utils"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/sd/lb"
	"github.com/openzipkin/zipkin-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				isSuccess := true
				if err != nil {
					isSuccess = false
				} else {
					switch rsp := response.(type) {
					case Response:
						if rsp.Error != nil {
							isSuccess = false
						}
					case *Response:
						if rsp != nil && rsp.Error != nil {
							isSuccess = false
						}
					default:
					}
				}
				duration.With("success", fmt.Sprintf("%v", isSuccess)).Observe(float64(time.Since(begin).Microseconds()))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(name string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			traceid := TraceID(ctx)

			defer func(begin time.Time) {
				var (
					reqstr string
					rspstr string
					rsperr error
				)

				switch request.(type) {
				case *http.Request:
					return
				default:
					reqstr = utils.Stringify(request)
				}

				switch rsp := response.(type) {
				case Response:
					rspstr = utils.Stringify(rsp.Value)
					rsperr = rsp.Error
				case *Response:
					rspstr = utils.Stringify(rsp.Value)
					rsperr = rsp.Error
				default:
					rspstr = utils.Stringify(rsp)
				}
				if err != nil || rsperr != nil {
					logger.Log(name, "transport_error", err, "err", rsperr, "took", time.Since(begin), "req", reqstr, "rsp", rspstr, "traceid", traceid)
				} else {
					logger.Log(name, "took", time.Since(begin), "req", reqstr, "rsp", rspstr, "traceid", traceid)
				}
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func TraceID(ctx context.Context) string {
	if parent := zipkin.SpanFromContext(ctx); parent != nil {
		return parent.Context().TraceID.String()
	}
	return ""
}

func TraceIDField(ctx context.Context) zap.Field {
	return zap.String("traceid", TraceID(ctx))
}

func makeLibError(err error) lib.Error {
	var (
		liberr lib.Error
		ok     bool
	)

	if retryerr, retry := err.(lb.RetryError); retry {
		liberr, ok = retryerr.Final.(lib.Error)
	} else {
		liberr, ok = err.(lib.Error)
	}

	if ok {
		return liberr
	} else {
		return lib.Error{Code: lib.ErrInternal, Message: err.Error()}
	}
}

func ErrorWrapperForServer(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		rsp, err := next(ctx, request)
		if err != nil {
			return rsp, err
		}

		// Server endpoint should return Response{}
		kitxrsp, ok := rsp.(Response)
		if !ok || kitxrsp.Error == nil {
			return rsp, nil
		}

		liberr := makeLibError(kitxrsp.Error)
		kitxrsp.Error = liberr

		// set the error via grpc metadata
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-error", fmt.Sprintf("%d,%s", liberr.Code, liberr.Message)))

		return kitxrsp, nil
	}
}

func ErrorWrapperForClient(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		rsp, err := next(ctx, request)
		if err == nil {
			return rsp, nil
		}
		err = makeLibError(err)
		return rsp, err
	}
}
