package util

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

func IotLog(formatStr string, a ...any) {
	slog.Info(fmt.Sprintf(formatStr, a...))
}

func IotLogErrWithStr(s string, err error) {
	AdditionalLog(err)
	slog.Error(s, err)
}

func IotLogError(err error) {
	AdditionalLog(err)
	if err != nil {
		slog.Error("", err)
	}
}

func IotLogErrorStr(s string) {
	AdditionalLogStr(s)
	slog.Error(s)
}

// func IotLogFatal(err error) {
// 	AdditionalLog(err)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func IotLogInfo(s string) {
	AdditionalLogStr(s)
	slog.Info(s)
}

type LogConfigParams struct {
	MaxSize    int `json:"max_size"`
	MaxBackups int `json:"max_backups"`
	MaxAge     int `json:"software_version"`
	LogLevel   int `json:"log_level"`
}

func ConfigLogFile(logFilename string, params LogConfigParams) {
	r := &lumberjack.Logger{
		Filename:   "./logs/" + logFilename,
		MaxSize:    params.MaxSize,    //日志最大的大小（M）
		MaxBackups: params.MaxBackups, //备份个数
		MaxAge:     params.MaxAge,     //最大保存天数（day）
		Compress:   false,             //是否压缩
		LocalTime:  false,
	}

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, r), nil))
	slog.SetDefault(logger)
}
