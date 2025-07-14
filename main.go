package main

import (
	"ayode.org/visor/client"
	"ayode.org/visor/config"
	"fmt"
	"log/slog"
	"os"
)

func incorrectArg() {
	fmt.Printf("That is not a valid subcommand!")
	os.Exit(1)
}

func runCfg() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if len(os.Args) < 3 {
		slog.Error("No configuration path provided!")
		os.Exit(2)
		return
	}
	cfgPath := os.Args[2]
	cfg := config.Parse(cfgPath)
	c := client.New(cfg)
	c.Execute()
}

func main() {
	if len(os.Args) < 2 {
		incorrectArg()
	}

	switch os.Args[1] {
	case "run":
		runCfg()
		break
	default:
		incorrectArg()
	}
}
