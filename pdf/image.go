package pdf

import (
	"context"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"runtime"
	
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/mandykoh/prism/srgb"	
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism"
)


/*

                        c := color.NRGBA{R: b[i], G: b[i+1], B: b[i+2], A: alpha}

                        ac, alpha2 := adobergb.ColorFromNRGBA(c)     // Interpret image pixel as Adobe RGB and convert to linear representation
                        sc := srgb.ColorFromXYZ(ac.ToXYZ())         // Convert to XYZ, then from XYZ to sRGB linear representation
                        c = sc.ToNRGBA(alpha2)

	br := bytes.NewReader(buf.Bytes())
	md, _, _:= autometa.Load(br)

	golog.Println(md)

	pr, err := md.ICCProfile()
	golog.Println(pr, err)

*/

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

	for idx, raw_im := range raw_images {

		im, _, err := image.Decode(raw_im)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode image, %w", err)
		}

		// https://pkg.go.dev/github.com/mandykoh/prism
		
		inputImg := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
		convertedImg := image.NewNRGBA(inputImg.Rect)
		
		for i := inputImg.Rect.Min.Y; i < inputImg.Rect.Max.Y; i++ {
			for j := inputImg.Rect.Min.X; j < inputImg.Rect.Max.X; j++ {
				inCol, alpha := adobergb.ColorFromNRGBA(inputImg.NRGBAAt(j, i))
				outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
				convertedImg.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
			}
		}
		
		images[idx] = convertedImg
		
	}
	
	return images, nil
}
