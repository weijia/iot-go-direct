//go:build linux

package util

import (
	"github.com/coreos/go-systemd/journal"
)

func AdditionalLog(err error) {
	err = journal.Print(journal.PriErr, "Error: ", err)
}

func AdditionalLogStr(s string) {
	err = journal.Print(journal.PriErr, "Error: ", s)
}
