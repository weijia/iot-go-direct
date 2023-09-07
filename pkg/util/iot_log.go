package util

import (
	"log"
)

func IotLogError(err error) {
	AdditionalLog(err)
	if err != nil {
		log.Error(err)
	}
}

func IotLogErrorStr(s string) {
	AdditionalLogStr(s)
	log.Error(s)
}

func IotLogFatal(err error) {
	AdditionalLog(err)
	if err != nil {
		log.Fatal(err)
	}
}
