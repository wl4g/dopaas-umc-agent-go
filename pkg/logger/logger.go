/**
 * Copyright 2017 ~ 2025 the original author or authors[983708408@qq.com].
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package logger

import (
	"github.com/wl4g/super-devops-umc-agent/pkg/config"
	"github.com/wl4g/super-devops-umc-agent/pkg/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

//
// Initialize zap
//

// See: pkg/logging_config.go
var Main *zapLoggerWrapper
var Receive *zapLoggerWrapper

// -------------------------
// Zap logger wrapper.
// -------------------------

// Customize logger wrapper.
type zapLoggerWrapper struct {
	*zap.Logger
	zapcore.Level
}

// Add method to zapLogger wrapper.
func (log zapLoggerWrapper) IsDebug() bool {
	return log.Level <= zapcore.DebugLevel
}

func (log zapLoggerWrapper) IsInfo() bool {
	return log.Level <= zapcore.InfoLevel
}

func (log zapLoggerWrapper) IsWarn() bool {
	return log.Level <= zapcore.WarnLevel
}

func (log zapLoggerWrapper) IsError() bool {
	return log.Level <= zapcore.ErrorLevel
}

func (log zapLoggerWrapper) IsFatal() bool {
	return log.Level <= zapcore.FatalLevel
}

//
// Create zapLogger wrapper.
//
func newZapLoggerWrapper(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, service string) *zapLoggerWrapper {
	return &zapLoggerWrapper{
		newZapLogger(filePath, level,
			maxSize, maxBackups, maxAge, compress, service),
		level,
	}
}

// -------------------------
// Zap logger creation.
// -------------------------

func InitZapLogger() {
	var logItems = config.GlobalConfig.Logging.LogItems

	// Init main logger.
	var mainLog = logItems[constant.DefaultLogMain]
	Main = newZapLoggerWrapper(
		mainLog.FileName,
		parseLogLevel(mainLog.Level),
		mainLog.Policy.MaxSize,
		mainLog.Policy.MaxBackups,
		mainLog.Policy.RetentionDays,
		true, constant.DefaultLogMain)

	// Init receive logger.
	var receiveLog = logItems[constant.DefaultLogReceive]
	Receive = newZapLoggerWrapper(
		receiveLog.FileName,
		parseLogLevel(receiveLog.Level),
		receiveLog.Policy.MaxSize,
		receiveLog.Policy.MaxBackups,
		receiveLog.Policy.RetentionDays,
		true, constant.DefaultLogReceive)
}

//
// [Create ZAP logger objects]
//
// filePath - logger file path
// level - logger level
// maxSize - Maximum size unit saved per logger file: M
// maxBackups - How many backups can logger files be saved at most
// maxAge - How many days can a file be saved at most?
// compress - Compression or not
// serviceName - service name
//
func newZapLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, service string) *zap.Logger {
	zapcoreObj := createZapCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(zapcoreObj, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("service", service)))
}

//
// [Create ZapCore objects]
//
func createZapCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	// 日志文件路径配置
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	// 公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     timeEncoder,                    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // (全/短)路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

// Parse zap-zapLog level
func parseLogLevel(text string) zapcore.Level {
	switch string(text) {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel
	case "panic", "PANIC":
		return zapcore.PanicLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Output time encoder
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("06-01-02 15:04:05"))
}
