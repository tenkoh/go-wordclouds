package wordclouds

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/fogleman/gg"
)

// Mask creates a slice of box structs from a given mask image to be passed to wordclouds.MaskBoxes.
func Mask[T ~string | ~[]byte](maskImage T, width int, height int, exclude color.RGBA) []*Box {
	res := make([]*Box, 0)

	var img image.Image
	switch v := any(maskImage).(type) {
	case string:
		i, err := gg.LoadPNG(v)
		if err != nil {
			panic(err)
		}
		img = i
	case []byte:
		i, _, err := image.Decode(bytes.NewReader(v))
		if err != nil {
			panic(err)
		}
		img = i
	}

	// scale
	imgw := img.Bounds().Dx()
	imgh := img.Bounds().Dy()

	wr := float64(width) / float64(imgw)
	wh := float64(height) / float64(imgh)
	scalingRatio := math.Min(wr, wh)
	// center
	xoffset := 0.0
	yoffset := 0.0
	if scalingRatio*float64(imgw) < float64(width) {
		xoffset = (float64(width) - scalingRatio*float64(imgw)) / 2
		res = append(res, &Box{
			float64(height),
			0.0,
			xoffset,
			0,
		})
		res = append(res, &Box{
			float64(height),
			float64(width) - xoffset,
			float64(width),
			0,
		})
	}

	if scalingRatio*float64(imgh) < float64(height) {
		yoffset = (float64(height) - scalingRatio*float64(imgh)) / 2
		res = append(res, &Box{
			yoffset,
			0.0,
			float64(width),
			0,
		})
		res = append(res, &Box{
			float64(height),
			0.0,
			float64(width),
			float64(height) - yoffset,
		})
	}
	step := 3
	bounds := img.Bounds()
	for i := bounds.Min.X; i < bounds.Max.X; i = i + step {
		for j := bounds.Min.Y; j < bounds.Max.Y; j = j + step {
			r, g, b, a := img.At(i, j).RGBA()
			er, eg, eb, ea := exclude.RGBA()

			if r == er && g == eg && b == eb && a == ea {
				b := &Box{
					math.Min(float64(j+step)*scalingRatio+yoffset, float64(height)),
					float64(i)*scalingRatio + xoffset,
					math.Min(float64(i+step)*scalingRatio+xoffset, float64(width)),
					float64(j)*scalingRatio + yoffset,
				}
				res = append(res, b)
			}
		}
	}

	return res
}
