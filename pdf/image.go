package pdf

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/aaronland/go-image/background"
	"github.com/aaronland/go-image/colour"
	"github.com/aaronland/go-image/imaging"
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
		new_im := colour.ToAdobeRGB(im)

		// Account for the fact that everything in PDF-land is upside down

		new_im = imaging.Rotate180(imaging.FlipV(new_im))
		new_im = imaging.Rotate180(new_im)

		// Draw image on white background

		final_im := background.AddBackground(new_im, backgroundColor)

		images[idx] = final_im

	}

	return images, nil
}
