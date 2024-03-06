package checkservice

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type CheckService struct {
	rootDir             string
	queue               chan string
	Logger              *slog.Logger
	minWidth, minHeight int
	files               []ImageInfo
}

func NewCheckService(rootDir string, maxThread int16, mv, mh int, logLevel slog.Level) *CheckService {
	return &CheckService{
		rootDir:   rootDir,
		queue:     make(chan string, maxThread),
		Logger:    slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})),
		minWidth:  mv,
		minHeight: mh,
		files:     []ImageInfo{},
	}
}

func (c *CheckService) ScanFiles() error {
	c.Logger.Info("scan files")
	err := filepath.Walk(c.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		c.Logger.Debug("scan file", "path", path)
		if !info.IsDir() {
			c.queue <- path
		}
		return nil
	})
	return err
}

func (c *CheckService) DoAction() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(c.queue)
		}()
		if err := c.ScanFiles(); err != nil {
			c.Logger.Error("error on scan files", "err", err)
			return
		}
	}()

	for p := range c.queue {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			imageConfig, err := c.CheckImage(path)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					c.Logger.Debug("file not image", "path", path, "err", err)
					return
				} else {
					c.Logger.Error("error on check image", "path", path, "err", err)
					return
				}
			}
			if !imageConfig.IsHighResolution {
				mu.Lock()
				c.files = append(c.files, imageConfig)
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()
}

func (c *CheckService) CheckImage(path string) (imageConfig ImageInfo, err error) {
	c.Logger.Debug("check file", "path", path)
	f, err := os.Open(path)
	if err != nil {
		return ImageInfo{}, nil
	}
	defer f.Close()
	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return ImageInfo{}, err
	}
	return NewImageInfo(config, path, c.minHeight, c.minWidth), nil
}

func (c *CheckService) DeleteFile() error {
	c.Logger.Info("delete small image files")
	for _, v := range c.files {
		c.Logger.Warn("delete file", "path", v.Path, "w", v.Config.Width, "h", v.Config.Height)
		if err := os.Remove(v.Path); err != nil {
			return err
		}
	}
	return nil
}

func (c *CheckService) PrintFile() error {
	c.Logger.Info("print small image files")
	for _, v := range c.files {
		c.Logger.Warn("small resolution file", "path", v.Path, "w", v.Config.Width, "h", v.Config.Height)
	}
	return nil
}
