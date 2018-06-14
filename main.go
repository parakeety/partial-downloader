package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/parakeety/partial-downloader/downloader"
)

var o = flag.String("o", "file", "filename to save downloaded file")

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		panic("please specify url")
	}

	d := &downloader.Downloader{
		URL:      flag.Args()[0],
		Filename: *o,
	}
	if err := d.Start(); err != nil {
		log.Fatalf("failed parallel download: %v", err)
	}
	fmt.Println("Parallel download success! Downloaded at: " + *o)
}
