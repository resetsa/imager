package checker

import (
	"image"
	"log/slog"
	"os"
)

type CheckImageResolution struct {
	logger        slog.Logger
	minResolution int
}

func (c *CheckImageResolution) Check(path string) (result bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return false, err
	}
	c.logger.Debug("check file", "path", path, "w", config.Width, "h", config.Height)
	result = min(config.Height, config.Width) <= c.minResolution
	return result, nil
}

func NewCheckImageResolution(minResolution int, logger slog.Logger) *CheckImageResolution {
	return &CheckImageResolution{minResolution: minResolution, logger: logger}
}
