package checkservice

import (
	"os"
	"path/filepath"
)

type Scanner interface {
	Scan(out chan string) error
}

type ScanDirectory struct {
	rootDir string
}

func NewScanDirectory(rootDir string) *ScanDirectory {
	return &ScanDirectory{rootDir: rootDir}
}

func (s *ScanDirectory) Scan(out chan string) error {
	return filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			out <- path
		}
		return nil
	})
}
