package zap

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Mode string

const (
	Production  Mode = "production"
	Development Mode = "development"
	None        Mode = "none"

	payloadKey = "payload"
)

var ErrUnsupportedZapLoggerMode = errors.New("unsupported zap mode")

func New(mode Mode) (l *zap.Logger, cleanup func(), err error) {
	var zapLogger *zap.Logger

	switch mode {
	case Development:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()
	case Production:
		config := zap.NewProductionConfig()
		zapLogger, err = config.Build()
	case None:
		zapLogger = zap.NewNop()
	default:
		err = fmt.Errorf("%w: %s", ErrUnsupportedZapLoggerMode, mode)
	}

	if err != nil {
		return nil, nil, err
	}

	undoRedirectStdLog := zap.RedirectStdLog(zapLogger)
	cleanup = func() {
		if errSync := zapLogger.Sync(); errSync != nil {
			if !strings.HasSuffix(errSync.Error(), "invalid argument") {
				log.Println(errSync)
			}
		}

		undoRedirectStdLog()
	}

	l = zapLogger.WithOptions(zap.AddCallerSkip(1)).With(zap.Namespace(payloadKey))

	return l, cleanup, nil
}
