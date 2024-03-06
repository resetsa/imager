package checkservice

import "image"

type ImageInfo struct {
	Config           image.Config
	Path             string
	IsHighResolution bool
}

func (i *ImageInfo) IsBetterRes(h, w int) bool {
	return min(i.Config.Height, i.Config.Width) >= min(h, w)
}

func NewImageInfo(config image.Config, path string, h, w int) ImageInfo {
	i := ImageInfo{
		Config: config,
		Path:   path,
	}
	i.IsHighResolution = i.IsBetterRes(h, w)
	return i
}
