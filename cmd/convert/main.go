package main

// pdf (preview app):  #E8BCE3
// png (preview app): #E1BEE2
// pdf (moz): #e8bde4
// png: #e8bde4

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"	
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"runtime"
	
	"github.com/aaronland/go-freeform/pdf"
	"github.com/aaronland/go-image-rotate/imaging"
	"github.com/sfomuseum/go-exif-update"
	_ "github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"	

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

		jpeg_opts := &jpeg.Options{
			Quality: 100,
		}
		
		for idx, im := range images {

			log.Printf("IM %T\n", im)
			log.Printf("MODEL %T\n", im.ColorModel())

			// https://pkg.go.dev/github.com/mandykoh/prism/adobergb
			tmp := image.NewNRGBA(im.Bounds())
			// tmp2 := image.NewNRGBA(im.Bounds())			
			
			srgb.EncodeImage(tmp, im, runtime.NumCPU())
			// adobergb.EncodeImage(tmp2, tmp, runtime.NumCPU())			

			im = tmp
			
			i := idx + 1

			jpeg_fname := strings.Replace(fname, ext, "", 1)
			jpeg_fname = fmt.Sprintf("%s-%03d.jpg", jpeg_fname, i)
			// jpeg_fname = fmt.Sprintf("%s-%03d.png", jpeg_fname, i)			
			jpeg_path := filepath.Join(root, jpeg_fname)

			temp_wr, err := os.CreateTemp("", "freeform.*.jpg")
			// temp_wr, err := os.CreateTemp("", "freeform.*.png")			

			if err != nil {
				log.Fatalf("Failed to create temp file for %s, %v", path, err)
			}

			defer os.Remove(temp_wr.Name())

			// Account for the fact that everything in PDF-land is upside down
			im = imaging.Rotate180(imaging.FlipV(im))
			im = imaging.Rotate180(im)

			
			dst := image.NewNRGBA(im.Bounds())

			backgroundColor := color.NRGBA{0xff, 0xff, 0xff, 0xff}
			
			draw.Draw(dst, dst.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)
			draw.Draw(dst, dst.Bounds(), im, im.Bounds().Min, draw.Over)


			// err = png.Encode(temp_wr, dst)

			err = jpeg.Encode(temp_wr, dst, jpeg_opts)

			if err != nil {
				log.Fatalf("Failed to write JPEG for %s, %v", jpeg_path, err)
			}

			err = temp_wr.Close()

			if err != nil {
				log.Fatalf("Failed to close %s, %v", jpeg_path, err)
			}

			os.Rename(temp_wr.Name(), jpeg_path)
			log.Println(temp_wr.Name())
			log.Println(jpeg_path)
			continue
			
			jpeg_r, err := os.Open(temp_wr.Name())

			if err != nil {
				log.Fatalf("Failed to open %s, %v", temp_wr.Name(), err)
			}

			defer jpeg_r.Close()

			jpeg_wr, err := os.OpenFile(jpeg_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatalf("Failed to open file for writing %s, %v", jpeg_path, err)
			}

			jpeg_dt := tm.Format(time.RFC3339)

			exif_props := map[string]interface{}{
				"DateTime":          jpeg_dt,
				"DateTimeDigitized": jpeg_dt,
				"DateTimeOriginal":  jpeg_dt,
			}

			err = update.UpdateExif(jpeg_r, jpeg_wr, exif_props)

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
