package checker

import "testing"

func TestCheckDirContent_Check(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		c          *CheckDirContent
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "not exist",
			c:    NewCheckDirContent(false, []string{".jpg"}, nil),
			args: args{
				path: "c:/temp/jpg/bird-of-paradise-3840x2160-15715.jpg",
			},
			wantResult: false,
			wantErr:    true,
		},
		{
			name: "check dir with subfolder deny",
			c:    NewCheckDirContent(false, []string{".jpg"}, nil),
			args: args{
				path: " C:/temp/jpg/",
			},
			wantResult: false,
			wantErr:    true,
		},
		{
			name: "check dir with subfolder permit",
			c:    NewCheckDirContent(true, []string{".jpg"}, nil),
			args: args{
				path: "C:/temp/jpg/",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "dir ok",
			c:    NewCheckDirContent(true, []string{".jpg", ".png"}, nil),
			args: args{
				path: "c:/temp/jpg/One-shot/[2014-09][JAP] Unknown Title (COMIC Kairakuten XTC Vol. 3)/",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "dir have wrong extension",
			c:    NewCheckDirContent(true, []string{".bmp"}, nil),
			args: args{
				path: "c:/temp/jpg/One-shot/[2014-09][JAP] Unknown Title (COMIC Kairakuten XTC Vol. 3)/",
			},
			wantResult: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.c.Check(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDirContent.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("CheckDirContent.Check() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
