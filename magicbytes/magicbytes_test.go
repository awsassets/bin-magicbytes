package magicbytes

import (
	"context"
	"path/filepath"
	"testing"
)

func TestSearch(t *testing.T) {
	//Arrange
	dir := filepath.Dir(saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg=="))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metas := []*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}

	omf := func(path, metaType string) bool { return true }

	type args struct {
		ctx       context.Context
		targetDir string
		metas     []*Meta
		onMatch   OnMatchFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Init", args: args{ctx: ctx, targetDir: dir, metas: metas, onMatch: omf}, wantErr: false},
		{name: "Nil Ctx", args: args{ctx: nil, targetDir: dir, metas: metas, onMatch: omf}, wantErr: true},
		{name: "Empty dir", args: args{ctx: ctx, targetDir: "", metas: metas, onMatch: omf}, wantErr: true},
		{name: "Empty dir", args: args{ctx: ctx, targetDir: dir, metas: nil, onMatch: omf}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Search(tt.args.ctx, tt.args.targetDir, tt.args.metas, tt.args.onMatch); (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
