package actor

import (
	"archive/zip"
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	CbzExt = ".cbz"
)

type CreateCBZ struct {
	Logger       *slog.Logger
	SaveOriginal bool
	DestDir      string
}

func NewCBZAct(saveOrig bool, destDir string, logger *slog.Logger) *CreateCBZ {
	return &CreateCBZ{SaveOriginal: saveOrig, Logger: logger, DestDir: destDir}
}

func (c *CreateCBZ) ActMany(paths []string) error {
	for _, v := range paths {
		if err := c.ActOnce(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *CreateCBZ) ActOnce(path string) error {
	rootArchiveDir := os.DirFS(path)
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	c.Logger.Debug("add files from", "path", path)
	if err := w.AddFS(rootArchiveDir); err != nil {
		return err
	}
	// Make sure to check the error on Close.
	if err := w.Close(); err != nil {
		return err
	}
	// Open the file for writing.
	c.Logger.Info("save as", "path", filepath.Base(path)+CbzExt)
	f, err := os.Create(c.DestDir + "/" + filepath.Base(path) + CbzExt)
	if err != nil {
		return err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if !c.SaveOriginal {
		c.Logger.Debug("remove original", "path", path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}
