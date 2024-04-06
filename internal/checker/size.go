package checker

import (
	"errors"
	"log/slog"
	"os"
)

type CheckImageSize struct {
	logger  *slog.Logger
	minSize int64
}

func (c *CheckImageSize) CheckImage(path string) (result bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, errors.New("path is dir")
	}
	c.logger.Debug("check file", "path", path, "size", fileInfo.Size())
	result = fileInfo.Size() >= c.minSize
	c.logger.Debug("result", "path", path, "result", result)
	return result, nil
}

func NewCheckImageSize(minSize int64, logger *slog.Logger) *CheckImageSize {
	return &CheckImageSize{minSize: minSize, logger: logger}
}
