package checkservice

import (
	"image"
	"testing"
)

func TestImageInfo_IsBetterRes(t *testing.T) {
	type args struct {
		h int
		w int
	}
	tests := []struct {
		name string
		i    *ImageInfo
		args args
		want bool
	}{
		{
			name: "test1",
			i: &ImageInfo{
				Config: image.Config{
					Width:  100,
					Height: 100,
				},
			},
			args: args{
				h: 200,
				w: 200,
			},
			want: false,
		},
		{
			name: "test2",
			i: &ImageInfo{
				Config: image.Config{
					Width:  1000,
					Height: 1000,
				},
			},
			args: args{
				h: 200,
				w: 200,
			},
			want: true,
		},
		{
			name: "test3",
			i: &ImageInfo{
				Config: image.Config{
					Width:  1000,
					Height: 100,
				},
			},
			args: args{
				h: 200,
				w: 200,
			},
			want: false,
		},
		{
			name: "test4",
			i: &ImageInfo{
				Config: image.Config{
					Width:  200,
					Height: 200,
				},
			},
			args: args{
				h: 200,
				w: 200,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsBetterRes(tt.args.h, tt.args.w); got != tt.want {
				t.Errorf("ImageInfo.CmpResolution() = %v, want %v", got, tt.want)
			}
		})
	}
}
