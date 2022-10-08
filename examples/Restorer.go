package main

/*
#cgo LDFLAGS: -lmmdeploy
#include <stdlib.h>
#include "mmdeploy/restorer.h"
int _mmdeploy_restorer_apply(
mmdeploy_restorer_t restorer,
uint8_t* data,
int height,
int width,
mmdeploy_mat_t** res) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_restorer_apply(restorer, &mat, 1, res);
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
	imagePath := os.Args[3]
	outputPath := "images/output_images/OutputRestorer.png"
	var restorer C.mmdeploy_restorer_t = nil
	status := C.mmdeploy_restorer_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &restorer)
	defer C.mmdeploy_restorer_destroy(restorer)
	img, err := packages.LoadImageMM(imagePath)
	if err != nil {
		fmt.Println(err)
	}
	var restorerResult *C.mmdeploy_mat_t = nil

	status = C._mmdeploy_restorer_apply(
		restorer,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&restorerResult)
	restorerImg := packages.CoverMMimage2RBGAMM(
		unsafe.Pointer(restorerResult.data),
		int(restorerResult.height),
		int(restorerResult.width),
	)
	context := gg.NewContextForRGBA(&restorerImg)
	context.SavePNG(outputPath)
	fmt.Printf("return code: %v\n", int(status))
	defer C.mmdeploy_restorer_release_result((*C.struct_mmdeploy_mat_t)(restorerResult), C.int(1))
}
