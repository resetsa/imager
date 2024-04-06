package actor

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ConvertAct struct {
	Logger        *slog.Logger
	SaveOriginal  bool
	PrefixNewFile string
}

func (c *ConvertAct) ActMany(paths []string) error {
	c.Logger.Debug("count file for processing", "len", len(paths))
	for _, v := range paths {
		if err := c.ActOnce(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *ConvertAct) ActOnce(path string) error {
	tmpLogger := c.Logger.With("file", path)
	tmpLogger.Info("convert file")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	oFileInfo, _ := f.Stat()
	oSize := oFileInfo.Size()
	tmpLogger.Debug("get original size", "size", oSize)
	imgConfig, format, err := image.DecodeConfig(f)
	tmpLogger.Debug("get original image prop", "format", format, "w", imgConfig.Width, "h", imgConfig.Height)
	if err != nil {
		return err
	}
	// set on begin file
	f.Seek(0, 0)
	oImage, _, err := image.Decode(f)
	if err != nil {
		tmpLogger.Error("error on decode", "err", err)
		return err
	}
	tmpLogger.Debug("resizing file")
	nImage := imaging.Resize(oImage, imgConfig.Width, imgConfig.Height, imaging.Lanczos)
	if err = c.Save(nImage, path, oSize); err != nil {
		return err
	}
	return nil
}

func (c *ConvertAct) Save(img *image.NRGBA, path string, sourceSize int64) error {
	tmpLogger := c.Logger.With("file", path)
	nPath := c.newPathGenerate(path)
	format, _ := imaging.FormatFromFilename(nPath)
	buf := bytes.NewBuffer([]byte{})
	if err := imaging.Encode(buf, img, format); err != nil {
		return err
	}
	if int64(buf.Len()) >= sourceSize {
		tmpLogger.Warn("original file smaller than convert", "original size", sourceSize, "new size", buf.Len())
		return nil
	}
	fo, err := os.Create(nPath)
	if err != nil {
		return err
	}
	defer fo.Close()
	tmpLogger.Info("save to", "new path", nPath, "size", buf.Len())
	if _, err = buf.WriteTo(fo); err != nil {
		return err
	}

	return nil
}

func (c *ConvertAct) newPathGenerate(path string) string {
	if !c.SaveOriginal {
		return path
	}
	fileName, dirName, ext := filepath.Base(path), filepath.Dir(path), filepath.Ext(path)
	nFileName := strings.ReplaceAll(fileName, ext, c.PrefixNewFile+ext)
	nPath := filepath.Join(dirName, nFileName)
	return nPath
}

func NewConvertAct(saveOrig bool, prefix string, logger *slog.Logger) *ConvertAct {
	if prefix == "" {
		return &ConvertAct{SaveOriginal: saveOrig, PrefixNewFile: "_mod", Logger: logger}
	}
	return &ConvertAct{SaveOriginal: saveOrig, PrefixNewFile: prefix, Logger: logger}
}
