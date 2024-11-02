package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	Log *logger
)

type logger struct {
	*zap.SugaredLogger
}

// New is a function to initialize logger
/*
 * debug bool - is debug mode
 * timeZone string - logger time zone, by default "GMT"
 */
func New(debug bool, timeZone string) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Цветная подсветка уровней
		EncodeTime:     customTimeEncoder,                // Кастомный формат времени
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Краткий формат caller
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	if timeZone != "" {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(time.FixedZone(timeZone, 3*60*60)).Format("2006-01-02 15:04:05"))
		}
	}

	var level zapcore.Level
	if debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	log := zap.New(core, zap.AddCaller())

	Log = &logger{
		SugaredLogger: log.Sugar(),
	}
}

// customTimeEncoder форматирует время в GMT+0
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02 15:04:05"))
}
