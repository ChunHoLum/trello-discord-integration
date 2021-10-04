package main

import (
	"context"
	"flag"
	"os"

	logger "github.com/ChunHoLum/trello-discord-bot/lib/logger"
)

const (
	configFilePath = "/etc/tdi.toml"
)

func main() {
	logger.Init()
	// command line argument parsing
	path := flag.String("config", configFilePath, "Trello discord integration config file path")
	flag.Parse()

	if err := run(*path); err != nil {
		// Exits with nonzero exit code and prints an error to a log.
		logger.Standard().WithError(err).Errorf("Terminating with fatal error...")
		os.Exit(1)
	} else {
		logger.Standard().Info("Successfully shut down")
	}
}

func run(configPath string) error {

	conf, err := LoadConfig(configPath)

	if err != nil {
		return err
	}

	logConfig := conf.Log

	if err = logger.Setup(logConfig); err != nil {
		return err
	}

	app, err := NewApp(*conf)

	if err != nil {
		return err
	}

	app.Run(context.Background())

	return nil
}
