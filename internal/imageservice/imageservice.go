package imageservice

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"resetsa/imager/internal/actor"
	"resetsa/imager/internal/checker"
	"resetsa/imager/internal/scanner"
	"sync"
)

type ImageService struct {
	scanner      scanner.Scanner
	queueForScan chan string
	queueForAct  chan string
	Logger       *slog.Logger
	filesForAct  []string
	actor        actor.Actor
	checker      checker.Checker
}

func NewImageService(rootDir string, maxThread int16, logger *slog.Logger, actor actor.Actor, checker checker.Checker) *ImageService {
	return &ImageService{
		scanner:      scanner.NewScanDirectory(rootDir),
		queueForScan: make(chan string, maxThread),
		queueForAct:  make(chan string, maxThread),
		Logger:       logger,
		filesForAct:  []string{},
		actor:        actor,
		checker:      checker,
	}
}

func (c *ImageService) ScanFiles() error {
	c.Logger.Info("scan files")
	return c.scanner.Scan(c.queueForScan)
}

func (c *ImageService) DoCheck() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(c.queueForScan)
		}()
		if err := c.ScanFiles(); err != nil {
			c.Logger.Error("error on scan files", "err", err)
			return
		}
	}()

	for p := range c.queueForScan {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			checkResult, err := c.checker.CheckImage(path)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					c.Logger.Debug("file not image", "path", path, "err", err)
					return
				} else {
					c.Logger.Error("error on check image", "path", path, "err", err)
					return
				}
			}
			if checkResult {
				mu.Lock()
				c.filesForAct = append(c.filesForAct, path)
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()
}

func (c *ImageService) Action() error {
	return c.actor.ActMany(c.filesForAct)
}

func (c *ImageService) DoAction() {
	defer close(c.queueForAct)
	var wg sync.WaitGroup
	for _, v := range c.filesForAct {
		wg.Add(1)
		c.queueForAct <- v
		go func(chan string) {
			defer wg.Done()
			imageConfig := <-c.queueForAct
			if err := c.actor.ActOnce(imageConfig); err != nil {
				c.Logger.Error("error on action", "err", err)
			}
		}(c.queueForAct)
	}
	wg.Wait()
}
