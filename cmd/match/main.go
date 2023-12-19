package main

import (
	"errors"
	"fmt"

	"github.com/TeruMiyake/match/client"
	"github.com/TeruMiyake/match/server"
	flag "github.com/spf13/pflag"
)

var (
	isServerMode bool
	isClientMode bool
)

func main() {
	initFlags()
	parseFlags()

	forceModeSelection()

	mode, err := getModeString()
	if err != nil {
		panic(err)
	}

	fmt.Println("Mode: ", mode)

	if isServerMode {
		server.RunServer()
	} else if isClientMode {
		client.RunClient()
	}
}

func initFlags() {
	flag.BoolVarP(&isServerMode, "server", "s", false, "server mode")
	flag.BoolVarP(&isClientMode, "client", "c", false, "client mode")
}

func parseFlags() {
	flag.Parse()
}

func forceModeSelection() {
	if isServerMode && isClientMode {
		fmt.Println("Error: cannot be both server and client")
		return
	}
	for !isServerMode && !isClientMode {
		promptForMode()
	}
}

func promptForMode() {
	fmt.Println("Select mode (s for server, c for client):")
	var mode string
	fmt.Scanln(&mode)
	if mode == "s" || mode == "S" {
		isServerMode = true
	} else if mode == "c" || mode == "C" {
		isClientMode = true
	} else {
		fmt.Println("Error: invalid string. (s for server, c for client)")
	}
}

func getModeString() (string, error) {
	if isServerMode {
		return "server", nil
	}
	if isClientMode {
		return "client", nil
	}
	return "", errors.New("no mode selected")
}
