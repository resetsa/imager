package checkservice

import (
	"image"
	"log/slog"
	"os"
)

type Checker interface {
	CheckImage(string) (ImageInfo, error)
}

type CheckerImageSize struct {
	logger  slog.Logger
	minSize int
}

func (c *CheckerImageSize) CheckImage(path string) (imageConfig ImageInfo, err error) {
	c.logger.Debug("check file", "path", path)
	f, err := os.Open(path)
	if err != nil {
		return ImageInfo{}, nil
	}
	defer f.Close()
	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return ImageInfo{}, err
	}
	return NewImageInfo(config, path, c.minSize), nil
}

func NewCheckerImageSize(minSize int, logger slog.Logger) *CheckerImageSize {
	return &CheckerImageSize{minSize: minSize, logger: logger}
}
