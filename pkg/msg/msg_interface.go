package msg

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	// "github.com/coreos/go-systemd/journal"
)

type BaseMsg struct {
	MsgType string `json:"msg_type"`
}

type MsgHandler interface {
	handle()
}

type MsgFactory interface {
	getTargetMsg() any
}

func HandleMsg(body []byte) {
	var base_msg BaseMsg
	if err := json.Unmarshal(body, &base_msg); err != nil {
		// err := journal.Print(journal.PriInfo, "This is an informational message.")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		if runtime.GOOS == "linux" {
			// Linux-specific code
			fmt.Println("Running on Linux")
			// err = journal.Print(journal.PriErr, "unmarshal error", err)
		} else if runtime.GOOS == "windows" {
			// Windows-specific code
			fmt.Println("Running on Windows")
		} else {
			// Code for other platforms
			fmt.Println("Running on another platform")
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	switch base_msg.MsgType {
	case "init":
		var init Init
		if err := json.Unmarshal(body, &init); err != nil {
			// err = journal.Print(journal.PriErr, "unmarshal error", err)
			if err != nil {
				log.Fatal(err)
			}
		}
		init.handle()
	}
}
