package util

import (
	"log"
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
