package main

/*
#cgo LDFLAGS: -lmmdeploy
#include <stdlib.h>
#include "mmdeploy/detector.h"
int _mmdeploy_detector_apply(mmdeploy_detector_t detector, uint8_t* data, int height, int width,
		mmdeploy_detection_t** res, int** res_count) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_detector_apply(detector, &mat, 1, res, res_count);
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

func CreateBBoxPointsMM(up float64, down float64, left float64, right float64) (points []packages.PointMM) {
	points = []packages.PointMM{
		packages.PointMM{left, up},
		packages.PointMM{right, up},
		packages.PointMM{right, down},
		packages.PointMM{left, down}}
	return points
}

func main() {
	deviceName := os.Args[1]
	modelPath := os.Args[2]
	imagePath := os.Args[3]
	outputPath := "/imgs/OutputObjectDetection.png"
	var detector C.mmdeploy_detector_t = nil
	var detection *C.mmdeploy_detection_t = nil
	var result_count *C.int = nil

	status := C.mmdeploy_detector_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &detector)
	defer C.mmdeploy_detector_destroy(detector)

	img, err := packages.LoadImageMM(imagePath)
	if err != nil {
		fmt.Println(err)
	}

	status = C._mmdeploy_detector_apply(
		detector,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&detection, &result_count)
	if int(status) == 0 {
		fmt.Println("It has been sucessfully applied the model.")
	}
	defer C.mmdeploy_detector_release_result(
		(*C.struct_mmdeploy_detection_t)(detection),
		result_count, 1)

	count := int(*result_count)
	dets := (*[1 << 28]C.mmdeploy_detection_t)(unsafe.Pointer(detection))[:count]
	context := gg.NewContextForRGBA(img)
	for _, det := range dets {
		bbox := &det.bbox
		fmt.Printf("left=%v, top=%v, right=%v, bottom=%v, label=%v, score=%v\n",
			bbox.left, bbox.top, bbox.right, bbox.bottom, det.label_id, det.score)
		points := CreateBBoxPointsMM(
			float64(bbox.top),
			float64(bbox.bottom),
			float64(bbox.left),
			float64(bbox.right))
		*context = packages.DrawBBoxMM(*context, points)
		context.SetRGBA(0, 1, 0, 1)
		context.SetLineWidth(3)
		context.Stroke()
	}
	context.SavePNG(outputPath)
}
