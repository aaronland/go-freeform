package main


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

		_, err = pdf.Images(ctx, r)

		if err != nil {
			log.Fatalf("Failed to read properties for %s, %v", path, err)
		}

	}

}
