package core

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.SugaredLogger

// debugLevel : debug info error warn fatal
// InitLogger init logger
func initLogger(lumberJackLogger lumberjack.Logger, debugLevel string) {
	writeSyncer := getLogWriter(lumberJackLogger)
	encoder := getEncoder()
	level := getLevel(debugLevel)

	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Filename: 日志文件的位置
// MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
// MaxBackups：保留旧文件的最大个数
// MaxAges：保留旧文件的最大天数
// Compress：是否压缩/归档旧文件
func getLogWriter(lumberJackLogger lumberjack.Logger) zapcore.WriteSyncer {
	// lumberJackLogger := &lumberjack.Logger{
	// 	Filename:   path,
	// 	MaxSize:    10,
	// 	MaxBackups: 5,
	// 	MaxAge:     30,
	// 	Compress:   false,
	// }
	// zapcore.AddSync(lumberJackLogger)
	// syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberJackLogger))
}

func getLevel(debugLevel string) zapcore.Level {
	level := zapcore.DebugLevel
	switch debugLevel {
	case "debug":
		level = zapcore.DebugLevel
		break
	case "info":
		level = zapcore.InfoLevel
		break
	case "error":
		level = zapcore.ErrorLevel
		break
	case "warn":
		level = zapcore.WarnLevel
		break
	case "fatal":
		level = zapcore.FatalLevel
		break
	}
	return level
}
