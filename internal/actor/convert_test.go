package actor

import (
	_ "image/jpeg"
	_ "image/png"
)

//func TestConvertAct_ActOnce(t *testing.T) {
//	type args struct {
//		path string
//	}
//	logLevel := slog.LevelDebug
//	tests := []struct {
//		name    string
//		c       *ConvertAct
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "test_jpg",
//			c:    NewConvertAct(true, "_mod", slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))),
//			args: args{
//				path: "c:/temp/jpg/bird-of-paradise-3840x2160-15715.jpg",
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := tt.c.ActOnce(tt.args.path); (err != nil) != tt.wantErr {
//				t.Errorf("ConvertAct.ActOnce() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func TestConvertAct_newPathGenerate(t *testing.T) {
//	type args struct {
//		path string
//	}
//	logLevel := slog.LevelDebug
//	tests := []struct {
//		name string
//		c    *ConvertAct
//		args args
//		want string
//	}{
//		{
//			name: "test1",
//			c:    NewConvertAct(true, "_mod", slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))),
//			args: args{
//				path: "c:/temp/jpg/bird-of-paradise-3840x2160-15715.jpg",
//			},
//			want: "c:\\temp\\jpg\\bird-of-paradise-3840x2160-15715_mod.jpg",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.newPathGenerate(tt.args.path); got != tt.want {
//				t.Errorf("ConvertAct.newPathGenerate() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
