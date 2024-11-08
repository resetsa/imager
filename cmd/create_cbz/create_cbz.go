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
	rootDir := pflag.StringP("source-dir", "s", "", "source dir")
	destDir := pflag.StringP("target-dir", "t", "", "dest dir")
	rmOrig := pflag.Bool("remove-orig", false, "remove original files")
	maxThreads := pflag.Int16P("goroutines", "g", 10, "max parallel process")
	debugLevel := pflag.BoolP("debug", "d", false, "debug level")
	pflag.Parse()

	if *rootDir == "" {
		usage()
		slog.Error("root_dir is required")
		os.Exit(1)
	}

	if *destDir == "" {
		usage()
		slog.Error("dest_dir is required")
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	if *debugLevel {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	actor := actor.NewCBZAct(!*rmOrig, *destDir, logger)
	scanner := scanner.NewScanDir(*rootDir)
	checker := checker.NewCheckDirContent(false, []string{".jpg", ".png", ".jpeg"}, logger)
	c := service.NewImageService(*rootDir, *maxThreads, logger, actor, checker, scanner)

	c.DoCheck()
	c.DoAction()

}
