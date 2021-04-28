package magicbytes

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func Test_searchMetasInAFile(t *testing.T) {
	//Arrange
	paths := []string{
		saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg=="),
		saveTestFile("/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////wgALCAABAAEBAREA/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxA="),
		saveTestFile("QVM="),
		saveTestFile("eC50eHQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAxMDA3NzcAMDAwMDAwMAAwMDAwMDAwADAwMDAwMDAwMDAyADE0MDQyMDY3NjcwADAwNzAyMgAgMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB1c3RhcgAwMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBUwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
	}
	ctx := context.Background()

	metas := []*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/jpeg", Offset: 0, Bytes: []byte{0xff, 0xd8, 0xff, 0xe0}},
		{Type: "application/x-tar", Offset: 0x101, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}},
		nil,
	}

	type args struct {
		ctx       context.Context
		inner_ctx context.Context
		path      string
		metas     *[]*Meta
	}
	tests := []struct {
		name    string
		args    args
		want    *Meta
		wantErr bool
	}{
		{name: "Valid png file.", args: args{ctx: ctx, inner_ctx: ctx, path: paths[0], metas: &metas}, want: metas[0], wantErr: false},
		{name: "Valid jpeg file.", args: args{ctx: ctx, inner_ctx: ctx, path: paths[1], metas: &metas}, want: metas[1], wantErr: false},
		{name: "InValid text file.", args: args{ctx: ctx, inner_ctx: ctx, path: paths[2], metas: &metas}, want: nil, wantErr: false},
		{name: "Invalid path", args: args{ctx: ctx, inner_ctx: ctx, path: "no-such-a-file.txt", metas: &metas}, want: nil, wantErr: true},
		{name: "Valid offset", args: args{ctx: ctx, inner_ctx: ctx, path: paths[3], metas: &metas}, want: metas[2], wantErr: false},
		{name: "InValid offset no return", args: args{ctx: ctx, inner_ctx: ctx, path: paths[3], metas: &[]*Meta{{Type: "application/x-tar", Offset: 0, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}}}}, want: nil, wantErr: false},
		{name: "Bigger offset value then the file size", args: args{ctx: ctx, inner_ctx: ctx, path: paths[3], metas: &[]*Meta{{Type: "application/x-tar", Offset: 123450, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}}}}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchMetasInAFile(tt.args.ctx, tt.args.inner_ctx, tt.args.path, tt.args.metas)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchMetasInAFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchMetasInAFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchMetasInAFile_CancelContext(t *testing.T) {
	//Arrange
	paths := []string{
		saveTestFile("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg=="),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	metas := []*Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
	}

	type args struct {
		ctx       context.Context
		inner_ctx context.Context
		path      string
		metas     *[]*Meta
	}
	tests := []struct {
		name    string
		args    args
		want    *Meta
		wantErr bool
	}{
		{name: "With cancelled context", args: args{ctx: ctx, inner_ctx: ctx, path: paths[0], metas: &metas}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchMetasInAFile(tt.args.ctx, tt.args.inner_ctx, tt.args.path, tt.args.metas)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchMetasInAFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchMetasInAFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func saveTestFile(b64 string) string {
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	dir, _ := ioutil.TempDir("", "")
	tmpfile, err := ioutil.TempFile(dir, "")
	if err != nil {
		log.Fatal(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(dec); err != nil {
		panic(err)
	}
	if err := tmpfile.Sync(); err != nil {
		panic(err)
	}

	return tmpfile.Name()
}
