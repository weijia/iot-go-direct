//go:build linux

package util

import (
	"log"

	"github.com/coreos/go-systemd/journal"
)

func IotLogFatal(err error) {
	err = journal.Print(journal.PriErr, "unmarshal error", err)
	if err != nil {
		log.Fatal(err)
	}
}
