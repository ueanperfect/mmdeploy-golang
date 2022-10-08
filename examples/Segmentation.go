package main

/*
#cgo LDFLAGS: -lmmdeploy
#include <stdlib.h>
#include "mmdeploy/segmentor.h"
int _mmdeploy_segmentation_apply(
mmdeploy_segmentor_t segmentor,
uint8_t* data,
int height,
int width,
mmdeploy_segmentation_t** res) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_segmentor_apply(segmentor, &mat, 1, res);
	return ec;
}
*/
import "C"
import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/ueanperfect/mmdeploy_golang/packages"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"
)

func main() {
	deviceName := os.Args[1]
	modelPath := os.Args[2]
	imgPath := os.Args[3]
	outputPath := "images/output_images/OutputSegmentation.png"
	var segmentor C.mmdeploy_segmentor_t = nil
	var segmentation *C.mmdeploy_segmentation_t = nil

	img, err := packages.LoadImageMM(imgPath)
	imgTemple, err1 := packages.LoadImageMM(imgPath)
	if err != nil {
		fmt.Println(err)
	}
	if err1 != nil {
		fmt.Println(err)
	}

	status := C.mmdeploy_segmentor_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &segmentor)
	defer C.mmdeploy_segmentor_destroy(segmentor)

	status = C._mmdeploy_segmentation_apply(
		segmentor,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&segmentation)

	PLATTE := packages.GeneratePlatteMM(int(segmentation.classes + 1))
	maskStartList := unsafe.Pointer(segmentation.mask)
	mask := packages.CoverListToMask(
		*imgTemple, maskStartList,
		int(segmentation.height),
		int(segmentation.width),
		PLATTE)
	output := packages.CombineImageMM(*img, mask, 0.5)
	content := gg.NewContextForRGBA(&output)
	content.SavePNG(outputPath)
	fmt.Printf("return code: %v\n", int(status))
	defer C.mmdeploy_segmentor_release_result((*C.struct_mmdeploy_segmentation_t)(segmentation), C.int(1))
}
