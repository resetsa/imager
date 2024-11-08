package main

import (
	"log/slog"
	"os"
	"resetsa/imager/internal/actor"
	"resetsa/imager/internal/checker"
	"resetsa/imager/internal/scanner"
	"resetsa/imager/internal/service"

	"github.com/spf13/pflag"
)

func Usage() {
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
		Usage()
		slog.Error("root_dir is required")
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	if *debugLevel {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	checker := checker.NewCheckImageResolution(*minSize, *logger)
	scanner := scanner.NewScanFile(*rootDir)

	var act actor.Actor
	act = &actor.PrintAct{Logger: logger}
	if *rmAction {
		act = &actor.DeleteAct{Logger: logger}
	}

	c := service.NewImageService(*rootDir, *maxThreads, logger, act, checker, scanner)

	c.DoCheck()
	c.DoAction()
}
