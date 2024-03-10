package checkservice

import "image"

type ImageInfo struct {
	Config           image.Config
	Path             string
	IsHighResolution bool
}

func (i *ImageInfo) IsBetterRes(size int) bool {
	return min(i.Config.Height, i.Config.Width) >= size
}

func NewImageInfo(config image.Config, path string, min int) ImageInfo {
	i := ImageInfo{
		Config: config,
		Path:   path,
	}
	i.IsHighResolution = i.IsBetterRes(min)
	return i
}
