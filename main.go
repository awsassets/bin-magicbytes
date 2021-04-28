package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/asalih/bin-magicbytes/magicbytes"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	m := []*magicbytes.Meta{
		{Type: "image/png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "image/jpeg", Offset: 0, Bytes: []byte{0xff, 0xd8, 0xff, 0xe0}},
		{Type: "application/x-tar", Offset: 0x101, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}},
		nil,
	}

	if err := magicbytes.Search(ctx, "C:\\tmp", m, func(path, metaType string) bool {
		fmt.Println(path)

		return false
	}); err != nil {
		log.Fatal(err)
	}

	//Add defer when removing the below
	cancel()

	fmt.Println("Waiting for input:")
	var input string
	fmt.Scanln(&input)
}
