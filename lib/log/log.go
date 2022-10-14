package log

import (
	"strings"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type options struct {
	level             string
	localTime         bool
	useCaller         bool
	disableStacktrace bool
}

var (
	std            logger
	defaultOptions = options{
		level:             "info",
		localTime:         false,
		useCaller:         true,
		disableStacktrace: false,
	}
)

func init() {
	_ = InitLogger()
}

type LoggerOption func(o *options)

func WithLevel(level string) LoggerOption {
	return func(o *options) {
		o.level = level
	}
}

func WithLocalTime(enable bool) LoggerOption {
	return func(o *options) {
		o.localTime = enable
	}
}

func WithCaller(enable bool) LoggerOption {
	return func(o *options) {
		o.useCaller = enable
	}
}

func DisableStacktrace(flag bool) LoggerOption {
	return func(o *options) {
		o.disableStacktrace = flag
	}
}

type Logger interface {
	Log(kvs ...interface{}) error
}

type logger struct {
	l   *zap.Logger
	s   *zap.SugaredLogger
	cfg *zap.Config
}

// impl kit.Logger
func (l *logger) Log(kv ...interface{}) error {
	if len(kv) > 0 && len(kv)%2 == 1 {
		msg, ok := kv[0].(string)
		if ok {
			l.s.Infow(msg, kv[1:]...)
			return nil
		}
	}

	l.s.Infow("", kv...)
	return nil
}

func GetLogger() kitlog.Logger {
	sugar := std.l.WithOptions(zap.AddCallerSkip(1)).Sugar()
	return &logger{
		l: sugar.Desugar(),
		s: sugar,
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	zapcore.ISO8601TimeEncoder(t, enc)
}

func InitLogger(opts ...LoggerOption) error {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	cfg := zap.NewProductionConfig()
	std.cfg = &cfg

	if o.localTime {
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg.EncoderConfig.EncodeTime = timeEncoder
	}

	if o.disableStacktrace {
		cfg.DisableStacktrace = true
	}

	cfg.DisableCaller = !o.useCaller
	SetLevel(o.level)

	if l, err := std.cfg.Build(); err != nil {
		return err
	} else {
		std.l = l
		std.s = l.Sugar()
	}

	return nil
}

func SetLevel(level string) {
	level = strings.ToLower(level)

	switch level {
	case "fatal":
		std.cfg.Level.SetLevel(zapcore.FatalLevel)
	case "panic":
		std.cfg.Level.SetLevel(zapcore.PanicLevel)
	case "error":
		std.cfg.Level.SetLevel(zapcore.ErrorLevel)
	case "warn":
		std.cfg.Level.SetLevel(zapcore.WarnLevel)
	case "info":
		std.cfg.Level.SetLevel(zapcore.InfoLevel)
	case "debug":
		std.cfg.Level.SetLevel(zapcore.DebugLevel)
	}
}

func L() *zap.Logger {
	return std.l
}

func S() *zap.SugaredLogger {
	return std.s
}
