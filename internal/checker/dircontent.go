package checker

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

type CheckDirContent struct {
	logger          *slog.Logger
	haveNestedDir   bool
	permitExtension []string
}

func (c *CheckDirContent) Check(path string) (result bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	entrys, err := f.ReadDir(0)
	if err != nil {
		return false, err
	}
	for _, entry := range entrys {
		if entry.IsDir() && !c.haveNestedDir {
			return false, errors.New("contains nested dir")
		}
		if !entry.IsDir() && !extInPermit(entry.Name(), c.permitExtension) {
			return false, errors.New("contains not permitted extension")
		}
	}
	return true, nil

}

func NewCheckDirContent(nestedDir bool, extensions []string, logger *slog.Logger) *CheckDirContent {
	return &CheckDirContent{haveNestedDir: nestedDir, permitExtension: extensions, logger: logger}
}

func extInPermit(path string, ext []string) bool {
	for _, e := range ext {
		if strings.Contains(path, e) {
			return true
		}
	}
	return false
}
