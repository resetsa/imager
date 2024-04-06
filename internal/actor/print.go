package actor

import "log/slog"

type PrintAct struct {
	Logger *slog.Logger
}

func (p *PrintAct) ActMany(paths []string) error {
	p.Logger.Debug("print small image files")
	for _, v := range paths {
		p.ActOnce(v)
	}
	return nil
}

func (p *PrintAct) ActOnce(path string) error {
	slog.Warn("small resolution file", "path", path)
	return nil
}
