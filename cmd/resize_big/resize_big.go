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

func usage() {
	pflag.PrintDefaults()
}

func main() {
	rootDir := pflag.StringP("root-dir", "r", "", "root dir for scan image files JPG/PNG")
	rmOrig := pflag.Bool("remove-orig", false, "remove original files")
	minSize := pflag.Int64P("size", "s", 50*1024, "min size for processing, bytes")
	maxResolution := pflag.Int64P("max-res", "m", 2048, "max size image in pixel")
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
	checker := checker.NewCheckImageSize(*minSize, logger)
	scanner := scanner.NewScanFile(*rootDir)

	act := actor.NewConvertAct(!*rmOrig, "", *maxResolution, logger)
	act.JpegQuality = 90
	c := service.NewImageService(*rootDir, *maxThreads, logger, act, checker, scanner)

	c.DoCheck()
	c.DoAction()
}
