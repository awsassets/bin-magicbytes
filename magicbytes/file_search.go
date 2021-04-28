package magicbytes

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
)

func searchMetasInAFile(ctx context.Context, path string, metas *[]*Meta) (*Meta, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	f, err := os.Open(path)
	if err != nil {
		//TODO: Add to log before return it. Make it in all returns
		log.Println(fmt.Errorf("Can't open file %s, error %v", path, err))

		return nil, err
	}

	defer func() { _ = f.Close() }()

	stat, err := f.Stat()
	if err != nil {
		log.Println(fmt.Errorf("Can't stat file %s, error %v", path, err))

		return nil, err
	}

	fileSize := stat.Size()

	for _, meta := range *metas {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if meta == nil {
			continue
		}

		//Offset and meta bytes shouldn't be bigger than file size.
		lenBytes := len(meta.Bytes)
		if meta.Offset+int64(lenBytes) > fileSize {
			//Seek doesn't retrun err for overflow seek and offset couldn't be bigger of the file size.
			continue
		}

		_, e := f.Seek(meta.Offset, 0)
		if e != nil {
			log.Println(fmt.Errorf("Can't seek on file %s, error %v", path, err))

			continue
		}

		mb := make([]byte, lenBytes)

		n, err := f.Read(mb)
		if err != nil {
			log.Println(fmt.Errorf("Can't read file %s, error %v", path, err))

			continue
		}

		if n == lenBytes && bytes.Equal(mb, meta.Bytes) {
			return meta, nil
		}
	}

	return nil, nil
}
