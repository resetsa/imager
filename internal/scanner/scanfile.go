package scanner

import (
	"os"
	"path/filepath"
)

type ScanFile struct {
	rootDir string
}

func NewScanFile(rootDir string) *ScanFile {
	return &ScanFile{rootDir: rootDir}
}

func (s *ScanFile) Scan(out chan string) error {
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

func (s *ScanFile) RootDir() string {
	return s.rootDir
}
