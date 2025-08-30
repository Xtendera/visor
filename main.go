package main

import (
	"fmt"
	"github.com/Xtendera/visor/client"
	"github.com/Xtendera/visor/config"
	"github.com/Xtendera/visor/util"
	"log/slog"
	"os"
)

func incorrectArg() {
	fmt.Printf("Invalid subcommand!\n")
	os.Exit(1)
}

func runCfg() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if len(os.Args) < 3 {
		slog.Error("No configuration path provided!\n")
		os.Exit(2)
		return
	}

	cfgPath := os.Args[2]
	cfg := config.Parse(cfgPath)
	c, err := client.New(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when initializing client: %s", err.Error()))
	}

	logger = logger.With("root", cfg.Root)
	slog.SetDefault(logger)

	c.Execute()
}

func printVersion() {
	fmt.Printf("Visor %s\n", util.GetVersion())
}

func main() {
	if len(os.Args) < 2 {
		incorrectArg()
	}

	switch os.Args[1] {
	case "run":
		runCfg()
		break
	case "version":
		printVersion()
		break
	default:
		incorrectArg()
	}
}
