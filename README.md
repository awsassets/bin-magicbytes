# bin-magicbytes
Package for detecting and verifying file type using magic bytes in pure Go
 
## Example Usage

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/asalih/bin-magicbytes/magicbytes"
)

func main() {
	var dirPath string
	flag.StringVar(&dirPath, "dirPath", ".", "Dir path")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := []*magicbytes.Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/jpeg", Offset: 0, Bytes: []byte{0xff, 0xd8, 0xff, 0xe0}},
		{Type: "application/x-tar", Offset: 0x101, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}},
		nil,
	}

	if err := magicbytes.Search(ctx, dirPath, m, func(path, metaType string) bool {
		fmt.Println(path)

		return false
	}); err != nil {
		log.Fatal(err)
	}
}
```

## API
`type Meta struct { ... }` Holds the name, magical bytes, and offset of the magical bytes to be searched.
`Search(...)` Search searches the given target directory to find files recursively using meta information.
`OnMatchFunc func(...)` OnMatchFunc represents a function to be called when Search function finds a match. Returning false must immediately stop Search process. For every match, onMatch callback is called concurrently.

