package main

import (
	"log/slog"
	"os"
	"resetsa/imager/internal/checkservice"

	"github.com/spf13/pflag"
)

func usage() {
	pflag.PrintDefaults()
}

func main() {
	rootDir := pflag.StringP("root-dir", "r", "", "root dir for scan image files JPG/PNG")
	rmAction := pflag.Bool("remove", false, "remove small image files")
	minSize := pflag.IntP("min", "m", 1024, "min side image")
	maxThreads := pflag.Int16P("threads", "t", 10, "max parallel process files")
	debugLevel := pflag.BoolP("debug", "d", false, "debug level")
	pflag.Parse()

	if *rootDir == "" {
		usage()
		slog.Error("root_dir is required")
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	if *debugLevel {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	checker := checkservice.NewCheckerImageSize(*minSize, *logger)

	var actor checkservice.Actor
	actor = &checkservice.PrintActor{Logger: logger}
	if *rmAction {
		actor = &checkservice.DeleteActor{Logger: logger}
	}
	c := checkservice.NewCheckService(*rootDir, *maxThreads, logger, actor, checker)

	c.DoCheck()

	if err := c.Action(); err != nil {
		c.Logger.Error("error on run action", "err", err)
	}
}
