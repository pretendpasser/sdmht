package entity

import (
	"context"

	"sdmht/lib/kitx"
	"sdmht/lib/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type internalTrace struct{}

func InjectInternalTrace(ctx context.Context, trace string) context.Context {
	if trace == "" {
		trace = uuid.New().String()
	}
	return context.WithValue(ctx, internalTrace{}, trace)
}

func SLog(ctx context.Context) *zap.SugaredLogger {
	trace := kitx.TraceID(ctx)
	if trace != "" {
		return log.L().With(zap.String("traceid", trace)).Sugar()
	}

	trace, _ = ctx.Value(internalTrace{}).(string)
	return log.L().With(zap.String("inner-traceid", trace)).Sugar()
}

func TraceId(ctx context.Context) string {
	trace := kitx.TraceID(ctx)
	if trace != "" {
		return trace
	}

	trace, _ = ctx.Value(internalTrace{}).(string)
	return trace
}
