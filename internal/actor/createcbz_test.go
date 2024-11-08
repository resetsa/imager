package actor

import (
	"log/slog"
	"os"
	"testing"
)

func TestCreateCBZ_ActOnce(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		c       *CreateCBZ
		args    args
		wantErr bool
	}{
		{
			name: "with_save_source",
			c:    NewCBZAct(true, "c:/temp", slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
			args: args{
				path: "C:/temp/jpg/One-shot/[2014-09][JAP] Unknown Title (COMIC Kairakuten XTC Vol. 3)",
			},
			wantErr: false,
		},
		//		{
		//			name: "without_save_source",
		//			c:    NewCreateCBZAct(false, "c:/temp", slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		//			args: args{
		//				path: "C:/temp/jpg/One-shot/[2014-10][JAP] Unknown Title (COMIC Kairakuten XTC Vol. 4)",
		//			},
		//			wantErr: false,
		//		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ActOnce(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateCBZ.ActOnce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
