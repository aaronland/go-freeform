package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-freeform/pdf"
	"github.com/aaronland/go-image/exif"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, path := range flag.Args() {

		root := filepath.Dir(path)
		fname := filepath.Base(path)
		ext := filepath.Ext(path)

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		props, err := pdf.Info(ctx, r)

		if err != nil {
			log.Fatalf("Failed to read properties for %s, %v", path, err)
		}

		tm, err := pdf.ParseDate(props["Creation date"])

		if err != nil {
			log.Fatalf("Failed to parse date for %s, %v", path, err)
		}

		_, err = r.Seek(0, 0)

		if err != nil {
			log.Fatalf("Failed to reset reader for %s, %v", path, err)
		}

		images, err := pdf.Images(ctx, r)

		if err != nil {
			log.Fatalf("Failed to derive images for %s, %v", path, err)
		}

		for idx, im := range images {

			i := idx + 1

			jpeg_fname := strings.Replace(fname, ext, "", 1)
			jpeg_fname = fmt.Sprintf("%s-%03d.jpg", jpeg_fname, i)
			jpeg_path := filepath.Join(root, jpeg_fname)

			jpeg_wr, err := os.OpenFile(jpeg_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatalf("Failed to open file for writing %s, %v", jpeg_path, err)
			}

			// https://github.com/rwcarlsen/goexif/blob/go1/exif/exif.go#L385
			jpeg_dt := tm.Format("2006:01:02 15:04:05")

			exif_props := map[string]interface{}{
				"DateTime":          jpeg_dt,
				"DateTimeDigitized": jpeg_dt,
				"DateTimeOriginal":  jpeg_dt,
				"Software":          "freeform",
			}

			err = exif.UpdateExif(im, jpeg_wr, exif_props)

			if err != nil {
				log.Fatalf("Failed to update EXIF data for %s, %w", jpeg_path, err)
			}

			err = jpeg_wr.Close()

			if err != nil {
				log.Fatalf("Failed to close %s, %v", jpeg_path, err)
			}

			log.Printf("Wrote %s\n", jpeg_path)
		}

	}

}
