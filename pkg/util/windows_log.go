//go:build windows

package util

import "log"

func IotLog(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
