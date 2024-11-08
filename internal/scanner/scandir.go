package scanner

import (
	"os"
	"path/filepath"
)

type ScanDir struct {
	rootDir string
}

func NewScanDir(rootDir string) *ScanDir {
	return &ScanDir{rootDir: rootDir}
}

func (s *ScanDir) Scan(out chan string) error {
	return filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			out <- path
		}
		return nil
	})
}

func (s *ScanDir) RootDir() string {
	return s.rootDir
}
