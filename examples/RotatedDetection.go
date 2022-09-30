package main

/*
#cgo CFLAGS: -I/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/install/include
#cgo LDFLAGS: -L/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#cgo LDFLAGS: -lmmdeploy -Wl,-rpath=/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#include <stdlib.h>
#include "mmdeploy/detector.h"
#include "mmdeploy/rotated_detector.h"
int _mmdeploy_rotated_detector_apply(
mmdeploy_rotated_detector_t detector,
uint8_t* data,
int height,
int width,
mmdeploy_rotated_detection_t** res,
int** res_count) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_rotated_detector_apply(detector, &mat, 1, res, res_count);
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
	"math"
	"os"
	"unsafe"
)

func CovertRbbox2Points(xc float64, yc float64, w float64, h float64, angle float64) (points []packages.PointMM) {
	wx := w / 2 * math.Acos(angle)
	wy := w / 2 * math.Asin(angle)
	hx := -h / 2 * math.Asin(angle)
	hy := h / 2 * math.Acos(angle)
	p1 := packages.PointMM{xc - wx - hx, yc - wy - hy}
	p2 := packages.PointMM{xc + wx - hx, yc + wy - hy}
	p3 := packages.PointMM{xc + wx + hx, yc + wy + hy}
	p4 := packages.PointMM{xc - wx + hx, yc - wy + hy}
	points = append(points, p1)
	points = append(points, p2)
	points = append(points, p3)
	points = append(points, p4)
	return points
}

func main() {
	deviceName := os.Args[1]
	modelPath := os.Args[2]
	imagePath := os.Args[3]
	var roated_detector C.mmdeploy_rotated_detector_t = nil
	status := C.mmdeploy_rotated_detector_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &roated_detector)
	defer C.mmdeploy_rotated_detector_destroy(roated_detector)
	img, err := packages.LoadImageMM(imagePath)
	if err != nil {
		fmt.Println(err)
	}
	var rotated_detection *C.mmdeploy_rotated_detection_t = nil
	var result_count *C.int = nil

	status = C._mmdeploy_rotated_detector_apply(
		roated_detector,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&rotated_detection,
		&result_count)

	fmt.Printf("return code: %v\n", int(status))
	count := int(*result_count)
	dets := (*[1 << 28]C.mmdeploy_rotated_detection_t)(unsafe.Pointer(rotated_detection))[:count]
	context := gg.NewContextForRGBA(img)
	for _, det := range dets {
		points := CovertRbbox2Points(
			float64(det.rbbox[0]),
			float64(det.rbbox[1]),
			float64(det.rbbox[2]),
			float64(det.rbbox[3]),
			float64(det.rbbox[4]))
		*context = packages.DrawBBoxMM(*context, points)
		context.SetRGBA(0, 1, 0, 1)
		context.SetLineWidth(3)
		context.Stroke()
	}
	context.SavePNG("OutputRotatedDetection.png")
	fmt.Printf("result count: %v\n", count)
	defer C.mmdeploy_rotated_detector_release_result((*C.struct_mmdeploy_rotated_detection_t)(rotated_detection), result_count)
}
