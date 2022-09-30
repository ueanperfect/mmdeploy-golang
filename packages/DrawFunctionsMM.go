package packages

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math/rand"
	"os"
	"unsafe"
)

func CoverListToMask(imgtemple image.RGBA, start unsafe.Pointer, height int, width int, platte []singlePlatteMM) (mask image.RGBA) {
	mask = imgtemple
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			index_number := i*width + j
			mask_class_number := *(*uint8)(unsafe.Pointer(uintptr(start) + 4*uintptr(index_number)))
			r, g, b := platte[mask_class_number].r, platte[mask_class_number].g, platte[mask_class_number].b
			mask.Set(j, i, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return mask
}

func GeneratePlatteMM(numClass int) (platte []singlePlatteMM) {
	for i := 0; i < numClass; i++ {
		r, g, b := rand.Intn(255), rand.Intn(255), rand.Intn(255)
		singlePlatte := singlePlatteMM{r, g, b}
		platte = append(platte, singlePlatte)
	}
	return platte
}

func CombineImageMM(img1 image.RGBA, img2 image.RGBA, value float64) (output image.RGBA) {
	output = img2
	for i := 0; i < len(img1.Pix); i++ {
		output.Pix[i] = uint8(float64(img1.Pix[i])*value) + uint8(float64(img2.Pix[i])*(1-value))
	}
	return output
}

func LoadImageMM(filePath string) (*image.RGBA, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tmp, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	img := image.NewRGBA(tmp.Bounds())
	bounds := tmp.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b, a := tmp.At(j, i).RGBA()
			img.Set(j, i, color.RGBA{R: uint8(b >> 8), G: uint8(g >> 8), B: uint8(r >> 8), A: uint8(a >> 8)})
		}
	}
	return img, err
}

func DrawBBoxMM(context gg.Context, points []PointMM) (context1 gg.Context) {
	for i := 0; i < len(points); i++ {
		context.LineTo(points[i].X, points[i].Y)
	}
	context.LineTo(points[0].X, points[0].Y)
	context1 = context
	return context1
}

func CoverMMimage2RBGAMM(startPoint unsafe.Pointer, height int, width int) (output image.RGBA) {
	point1 := image.Point{0, 0}
	point2 := image.Point{width, height}
	rec := image.Rectangle{point1, point2}
	output = *image.NewRGBA(rec)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			index_number := i*width + j
			r := *(*uint8)(unsafe.Pointer(uintptr(startPoint) + 3*uintptr(index_number)))
			b := *(*uint8)(unsafe.Pointer(uintptr(startPoint) + 3*uintptr(index_number+1)))
			g := *(*uint8)(unsafe.Pointer(uintptr(startPoint) + 3*uintptr(index_number+2)))
			output.Set(j, i, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return output
}
