package main

/*
#cgo CFLAGS: -I/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/install/include
#cgo LDFLAGS: -L/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#cgo LDFLAGS: -lmmdeploy -Wl,-rpath=/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#include <stdlib.h>
#include "mmdeploy/classifier.h"
int _mmdeploy_classification_apply(
mmdeploy_classifier_t classifier,
uint8_t* data,
int height,
int width,
mmdeploy_classification_t** res,
int** res_count) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_classifier_apply(classifier, &mat, 1, res, res_count);
	return ec;
}
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"
)

func loadImageClassification(filePath string) (*image.RGBA, error) {
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
	fmt.Printf("%v %v %v\n", height, width, img.Stride)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b, _ := tmp.At(j, i).RGBA()
			img.Set(j, i, color.RGBA{R: uint8(b >> 8), G: uint8(g >> 8), B: uint8(r >> 8), A: 255})
		}
	}
	return img, err
}
func main() {
	deviceName := os.Args[1]
	modelPath := os.Args[2]
	imagePath := os.Args[3]
	var classifier C.mmdeploy_classifier_t = nil
	var classification *C.mmdeploy_classification_t = nil
	var resultCount *C.int = nil

	status := C.mmdeploy_classifier_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &classifier)
	defer C.mmdeploy_classifier_destroy(classifier)
	img, err := loadImageClassification(imagePath)
	if err != nil {
		fmt.Println(err)
	}

	status = C._mmdeploy_classification_apply(
		classifier,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&classification,
		&resultCount)
	if int(status) == 0 {
		fmt.Println("It has been sucessfully applied the model.")
	}
	defer C.mmdeploy_classifier_release_result(
		(*C.struct_mmdeploy_classification_t)(classification),
		resultCount,
		C.int(1))

	count := int(*resultCount)
	result := (*[1 << 28]C.mmdeploy_classification_t)(unsafe.Pointer(classification))[:count]
	for _, det := range result {
		if det.score < 0.01 {
			continue
		}
		fmt.Printf("label: %v  score: %v \n", det.label_id, det.score)
	}
}
