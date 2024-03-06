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
	minWidth := pflag.IntP("min-width", "w", 1024, "min width")
	minHeight := pflag.IntP("min-height", "h", 768, "min height")
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

	c := checkservice.NewCheckService(*rootDir, *maxThreads, *minWidth, *minHeight, logLevel)
	c.DoAction()

	runAction := c.PrintFile
	if *rmAction {
		runAction = c.DeleteFile
	}
	if err := runAction(); err != nil {
		c.Logger.Error("error on run action", "err", err)
	}
}
