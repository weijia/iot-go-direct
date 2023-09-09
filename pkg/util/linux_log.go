//go:build linux

package util

import (
	"github.com/coreos/go-systemd/journal"
)

func AdditionalLog(err error) {
	journal.Print(journal.PriErr, "Error: ", err)
}

func AdditionalLogStr(s string) {
	journal.Print(journal.PriErr, "Error: ", s)
}
