package actor

import (
	"log/slog"
	"os"
)

type DeleteAct struct {
	Logger *slog.Logger
}

func (d *DeleteAct) ActMany(paths []string) error {
	d.Logger.Debug("delete small image files")
	for _, v := range paths {
		return d.ActOnce(v)
	}
	return nil
}

func (d *DeleteAct) ActOnce(path string) error {
	slog.Warn("delete file", "path", path)
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
