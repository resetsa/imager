package service

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
	Logger       *slog.Logger
	pathsForAct  []string
	actor        actor.Actor
	checker      checker.Checker
	limiter      chan struct{}
}

func NewImageService(rootDir string, maxThread int16, logger *slog.Logger, actor actor.Actor, checker checker.Checker, scanner scanner.Scanner) *ImageService {
	return &ImageService{
		scanner:      scanner,
		queueForScan: make(chan string),
		Logger:       logger,
		pathsForAct:  []string{},
		actor:        actor,
		checker:      checker,
		limiter:      make(chan struct{}, maxThread),
	}
}

func (c *ImageService) Scan() error {
	c.Logger.Info("scan targets", "path", c.scanner.RootDir())
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
		if err := c.Scan(); err != nil {
			c.Logger.Error("error on scan", "err", err)
			return
		}
	}()

	for p := range c.queueForScan {
		wg.Add(1)
		c.limiter <- struct{}{}
		go func(path string) {
			defer func() {
				wg.Done()
				<-c.limiter
			}()
			checkResult, err := c.checker.Check(path)
			if err != nil {
				if errors.Is(err, image.ErrFormat) {
					c.Logger.Debug("file not image", "path", path, "err", err)
					return
				}
				c.Logger.Error("error on check", "path", path, "err", err)
				return
			}
			if checkResult {
				mu.Lock()
				c.pathsForAct = append(c.pathsForAct, path)
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()
}

func (c *ImageService) DoAction() {
	var wg sync.WaitGroup
	for _, v := range c.pathsForAct {
		wg.Add(1)
		c.limiter <- struct{}{}
		go func(path string) {
			defer func() {
				wg.Done()
				<-c.limiter
			}()
			if err := c.actor.ActOnce(path); err != nil {
				c.Logger.Error("error on action", "err", err)
			}
		}(v)
	}
	wg.Wait()
}
