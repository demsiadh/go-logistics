package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	LevelKey      = "level"
	NameKey       = "logger"
	StacktraceKey = "stacktrace"
	MessageKey    = "msg"
)

var (
	Log *zap.Logger
)

// NewProductionLogger 创建生产环境下的日志记录器
func NewProductionLogger() *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LevelKey:       LevelKey,
			NameKey:        NameKey,
			StacktraceKey:  StacktraceKey,
			MessageKey:     MessageKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     CustomTimeEncoder, // 使用自定义编码器
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	), zap.AddCaller())
}

// CustomTimeEncoder 自定义时间编码器（带毫秒和时区）
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000 CST"))
}

func initLog() {
	Log = NewProductionLogger()
}
