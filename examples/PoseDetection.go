package main

/*
#cgo LDFLAGS: -lmmdeploy
#include <stdlib.h>
#include "mmdeploy/pose_detector.h"
int _mmdeploy_pose_detector_apply(
mmdeploy_pose_detector_t detector,
uint8_t* data,
int height,
int width,
mmdeploy_pose_detection_t** res) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_pose_detector_apply(detector, &mat, 1, res);
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
	outputPath := "images/output_images/OutputPoseDetection.png"
	var poseDetector C.mmdeploy_pose_detector_t = nil
	var poseDetection *C.mmdeploy_pose_detection_t = nil

	status := C.mmdeploy_pose_detector_create_by_path(C.CString(modelPath), C.CString(deviceName), C.int(0), &poseDetector)
	defer C.mmdeploy_pose_detector_destroy(poseDetector)

	img, err := packages.LoadImageMM(imagePath)
	if err != nil {
		fmt.Println(err)
	}

	status = C._mmdeploy_pose_detector_apply(
		poseDetector,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&poseDetection)
	if int(status) == 0 {
		fmt.Println("It has been sucessfully applied the model.")
	}
	defer C.mmdeploy_pose_detector_release_result((*C.struct_mmdeploy_pose_detection_t)(poseDetection), C.int(1))

	context := gg.NewContextForRGBA(img)
	points := (*[1 << 28]C.mmdeploy_point_t)(unsafe.Pointer(poseDetection.point))[:int(poseDetection.length)]
	for i := 0; i < int(poseDetection.length); i++ {
		point := points[i]
		context.DrawCircle(float64(point.x), float64(point.y), 5)
	}
	context.SetRGBA(0, 1, 0, 1)
	context.Fill()
	context.Stroke()
	context.SavePNG(outputPath)
}
