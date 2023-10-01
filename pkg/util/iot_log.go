package util

import (
	"log"
	"github.com/natefinch/lumberjack"
)

func IotLogError(err error) {
	AdditionalLog(err)
	if err != nil {
		log.Println(err)
	}
}

func IotLogErrorStr(s string) {
	AdditionalLogStr(s)
	log.Println(s)
}

func IotLogFatal(err error) {
	AdditionalLog(err)
	if err != nil {
		log.Fatal(err)
	}
}

func IotLogInfo(s string) {
	AdditionalLogStr(s)
	log.Println(s)
}

func IotLog(formatStr string, a ...any){
	log.Printf(formatStr, a...)
}

type LogConfigParams struct {
	MaxSize   int `json:"max_size"`
	MaxBackups int `json:"max_backups"`
	MaxAge int `json:"software_version"`
	LogLevel int `json:"log_level"`
}

func ConfigLogFile(logFilename string, params LogConfigParams) {
	hook := lumberjack.Logger{
		Filename:   "./logs/" + logFilename,
		MaxSize:    params.MaxSize,    //日志最大的大小（M）
		MaxBackups: params.MaxBackups,   //备份个数
		MaxAge:     params.MaxAge,    //最大保存天数（day）
		Compress:   false, //是否压缩
		LocalTime:  false,
	}
	log.SetOutput(&hook)
}