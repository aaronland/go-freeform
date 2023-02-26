package main

/*

> go run cmd/info/main.go ~/Downloads/Untitled\ 41.pdf
2023/02/26 11:43:38 Content creator:
2023/02/26 11:43:38 Modification date: D:20230222192640+00'00'
2023/02/26 11:43:38 Hybrid: No
2023/02/26 11:43:38 Linearized: No
2023/02/26 11:43:38 Permissions: Full access
2023/02/26 11:43:38 Page size: 479.69 x 733.23 points
2023/02/26 11:43:38 Author:
2023/02/26 11:43:38 Subject:
2023/02/26 11:43:38 Using XRef streams: No
2023/02/26 11:43:38 Encrypted: No
2023/02/26 11:43:38 PDF version: 1.4
2023/02/26 11:43:38 Page count: 1
2023/02/26 11:43:38 Acroform: No
2023/02/26 11:43:38 PDF Producer: iOS Version 16.3.1 (Build 20D67) Quartz PDFContext
2023/02/26 11:43:38 Thumbnails: No
2023/02/26 11:43:38 Tagged: Yes
2023/02/26 11:43:38 Using object streams: No
2023/02/26 11:43:38 Watermarked: No
2023/02/26 11:43:38 Title:
2023/02/26 11:43:38 Creation date: D:20230222192640+00'00'
*/

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-freeform/pdf"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		props, err := pdf.Info(ctx, r)

		if err != nil {
			log.Fatalf("Failed to read properties for %s, %v", path, err)
		}

		for k, v := range props {
			log.Printf("%s: %s\n", k, v)
		}
	}

}
