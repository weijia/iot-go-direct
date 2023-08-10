//go:build windows

package util

import "log"

func IotLogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
