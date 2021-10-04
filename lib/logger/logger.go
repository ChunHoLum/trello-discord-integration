package log

import (
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type Config struct {
	Output   string `toml:"output"`
	Severity string `toml:"severity"`
}

type contextKey struct{}

// InitLogger sets up logger for a typical daemon scenario until configuration
// file is parsed
func Init() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	log.SetOutput(os.Stdout)
}

func Setup(conf Config) error {
	switch conf.Output {
	case "stderr", "error", "2":
		log.SetOutput(os.Stderr)
	case "", "stdout", "out", "1":
		log.SetOutput(os.Stdout)
	default:
		// assume it's a file path:
		logFile, err := os.Create(conf.Output)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
	}

	switch strings.ToLower(conf.Severity) {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "err", "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	return nil
}

func Get(ctx context.Context) log.FieldLogger {
	if logger, ok := ctx.Value(contextKey{}).(log.FieldLogger); ok && logger != nil {
		return logger
	}

	return Standard()
}

func Standard() log.FieldLogger {
	return log.StandardLogger()
}

func withLogger(ctx context.Context, logger log.FieldLogger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func WithField(ctx context.Context, key string, value interface{}) (context.Context, log.FieldLogger) {
	logger := Get(ctx).WithField(key, value)
	return withLogger(ctx, logger), logger
}

func WithFields(ctx context.Context, logFields Fields) (context.Context, log.FieldLogger) {
	logger := Get(ctx).WithFields(logFields)
	return withLogger(ctx, logger), logger
}
