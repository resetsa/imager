package checkservice

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"sync"
)

type CheckService struct {
	scanner Scanner
	queue   chan string
	Logger  *slog.Logger
	files   []ImageInfo
	actor   Actor
	checker Checker
}

func NewCheckService(rootDir string, maxThread int16, logger *slog.Logger, actor Actor, checker Checker) *CheckService {
	return &CheckService{
		scanner: NewScanDirectory(rootDir),
		queue:   make(chan string, maxThread),
		Logger:  logger,
		files:   []ImageInfo{},
		actor:   actor,
		checker: checker,
	}
}

func (c *CheckService) ScanFiles() error {
	c.Logger.Info("scan files")
	return c.scanner.Scan(c.queue)
}

func (c *CheckService) DoCheck() {
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
			imageConfig, err := c.checker.CheckImage(path)
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

func (c *CheckService) Action() error {
	return c.actor.Act(c.files)
}
