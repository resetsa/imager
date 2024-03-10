package checkservice

import (
	"log/slog"
	"os"
)

type Actor interface {
	Act([]ImageInfo) error
}

type PrintActor struct {
	Logger *slog.Logger
}

func (p *PrintActor) Act(images []ImageInfo) error {
	p.Logger.Debug("print small image files")
	for _, v := range images {
		slog.Warn("small resolution file", "path", v.Path, "w", v.Config.Width, "h", v.Config.Height)
	}
	return nil
}

type DeleteActor struct {
	Logger *slog.Logger
}

func (d *DeleteActor) Act(images []ImageInfo) error {
	d.Logger.Debug("delete small image files")
	for _, v := range images {
		slog.Warn("delete file", "path", v.Path, "w", v.Config.Width, "h", v.Config.Height)
		if err := os.Remove(v.Path); err != nil {
			return err
		}
	}
	return nil
}
