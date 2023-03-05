package pdf

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"runtime"

	"github.com/aaronland/go-image-rotate/imaging"
	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Images(ctx context.Context, r io.ReadSeeker) ([]image.Image, error) {

	pages := []string{}
	conf := &pdfcpu.Configuration{}

	// At the end of the line these images are being extracted here:
	// func renderDeviceRGBToPNG(im *PDFImage, resourceName string) (io.Reader, string, error) {
	// github.com/pdfcpu/pdfcpu/pkg/pdfcpu/writeImage.go

	raw_images, err := api.ExtractImagesRaw(r, pages, conf)

	if err != nil {
		return nil, fmt.Errorf("Failed to extract images, %w", err)
	}

	images := make([]image.Image, len(raw_images))

	backgroundColor := color.NRGBA{0xff, 0xff, 0xff, 0xff}

	for idx, raw_im := range raw_images {

		im, _, err := image.Decode(raw_im)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode image, %w", err)
		}

		// Freeform uses Adobe RGB 1998
		// https://pkg.go.dev/github.com/mandykoh/prism

		inputImg := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
		new_im := image.NewNRGBA(inputImg.Rect)

		for i := inputImg.Rect.Min.Y; i < inputImg.Rect.Max.Y; i++ {
			for j := inputImg.Rect.Min.X; j < inputImg.Rect.Max.X; j++ {
				inCol, alpha := adobergb.ColorFromNRGBA(inputImg.NRGBAAt(j, i))
				outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
				new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
			}
		}

		// Account for the fact that everything in PDF-land is upside down

		new_im = imaging.Rotate180(imaging.FlipV(new_im))
		new_im = imaging.Rotate180(new_im)

		// Draw image on white background

		final_im := image.NewNRGBA(new_im.Bounds())

		draw.Draw(final_im, final_im.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)
		draw.Draw(final_im, final_im.Bounds(), new_im, new_im.Bounds().Min, draw.Over)

		images[idx] = final_im

	}

	return images, nil
}
