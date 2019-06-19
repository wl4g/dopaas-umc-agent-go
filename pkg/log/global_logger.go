//
// Copyright 2017 ~ 2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"umc-agent/pkg/common"
)

//
// [创建ZAP日志对象]
//
// filePath - 日志文件路径
// level - 日志级别
// maxSize - 每个日志文件保存的最大尺寸 单位：M
// maxBackups - 日志文件最多保存多少个备份
// maxAge - 文件最多保存多少天
// compress - 是否压缩
// serviceName - 服务名
//
func newZapLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	zapcoreObj := createZapCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(zapcoreObj, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
}

//
// [创建ZapCore对象]
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
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

//
// 初始化创建zap(logger)
//

var MainLogger *zap.Logger
var GatewayLogger *zap.Logger
var RedisLogger *zap.Logger
var LmdbLogger *zap.Logger
var HttpLogger *zap.Logger

func init() {
	MainLogger = newZapLogger(common.LOG_MAIN_FILENAME, zapcore.InfoLevel, 128, 30, 7, true, "Main")
	GatewayLogger = newZapLogger(common.LOG_GATEWAY_FILENAME, zapcore.DebugLevel, 128, 30, 7, true, "Gateway")
	HttpLogger = newZapLogger(common.LOG_HTTP_FILENAME, zapcore.InfoLevel, 128, 30, 7, true, "http")
	RedisLogger = newZapLogger(common.LOG_REDIS_FILENAME, zapcore.InfoLevel, 128, 30, 7, true, "redis")
	LmdbLogger = newZapLogger(common.LOG_LMDB_FILENAME, zapcore.InfoLevel, 128, 30, 7, true, "lmdb")
}
