package magicbytes

import (
	"bytes"
	"context"
	"os"
)

func searchMetasInAFile(ctx context.Context, inner_ctx context.Context, path string, metas *[]*Meta) (*Meta, error) {
	if !checkContextIsAlive(ctx, inner_ctx) {
		return nil, &ContextCancelledError{}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()

	for _, meta := range *metas {
		if !checkContextIsAlive(ctx, inner_ctx) {
			return nil, &ContextCancelledError{}
		}

		if meta == nil {
			continue
		}

		//Offset and meta bytes shouldn't be bigger than file size.
		len_bytes := len(meta.Bytes)
		if meta.Offset+int64(len_bytes) > fileSize {
			//Seek doesn't retrun err for overflow seek and offset couldn't be bigger of the file size.
			continue
		}

		_, e := f.Seek(meta.Offset, 0)
		if e != nil {
			//TODO: Add log here
			continue
		}

		mb := make([]byte, len_bytes)
		f.Read(mb)

		if bytes.Compare(mb, meta.Bytes) == 0 {
			return meta, nil
		}
	}

	return nil, nil
}
