package main

/*
#cgo CFLAGS: -I/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/install/include
#cgo LDFLAGS: -L/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#cgo LDFLAGS: -lmmdeploy -Wl,-rpath=/home/liyueyan/文档/mmdeploy_Project/mmdeploy/build/lib
#include <stdlib.h>
#include "mmdeploy/text_detector.h"
#include "mmdeploy/text_recognizer.h"
#include "mmdeploy/common.h"
int _mmdeploy_text_detector_apply(
mmdeploy_text_detector_t text_detector,
uint8_t* data,
int height,
int width,
mmdeploy_text_detection_t** res_bbox,
int** res_count) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_text_detector_apply(text_detector, &mat, 1, res_bbox, res_count);
	return ec;
}

int _mmdeploy_text_recognizer_apply_bbox(
mmdeploy_text_recognizer_t text_recognizer,
uint8_t* data,
int height,
int width,
mmdeploy_text_detection_t* bboxes,
int* bboxes_count,
mmdeploy_text_recognition_t** res_text) {
	mmdeploy_mat_t mat = {data, height, width, 4, MMDEPLOY_PIXEL_FORMAT_BGRA, MMDEPLOY_DATA_TYPE_UINT8};
	int ec;
	ec = mmdeploy_text_recognizer_apply_bbox(text_recognizer, &mat, 1, bboxes,bboxes_count,res_text);
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

func CovertBboxXY2Points(bbox [4]C.mmdeploy_point_t) (points []packages.PointMM) {
	for i := 0; i < len(bbox); i++ {
		fmt.Printf("point_%v: (x:%v y:%v) \n", i+1, float64(bbox[i].x), float64(bbox[i].y))
		points = append(points,
			packages.PointMM{float64(bbox[i].x), float64(bbox[i].y)},
		)
	}
	return points
}
func main() {
	deviceName := os.Args[1]
	detModelPath := os.Args[2]
	regModelPath := os.Args[3]
	imgPath := os.Args[4]
	outputPath := "/imgs/OutputOcr.png"

	var textDetector C.mmdeploy_text_detector_t = nil
	var textRecognizer C.mmdeploy_text_recognizer_t = nil
	var bboxes *C.mmdeploy_text_detection_t = nil
	var texts *C.mmdeploy_text_recognition_t = nil
	var resultCountBbox *C.int = nil

	status1 := C.mmdeploy_text_detector_create_by_path(C.CString(detModelPath), C.CString(deviceName), C.int(0), &textDetector)
	defer C.mmdeploy_text_detector_destroy(textDetector)

	status2 := C.mmdeploy_text_recognizer_create_by_path(C.CString(regModelPath), C.CString(deviceName), C.int(0), &textRecognizer)
	defer C.mmdeploy_text_recognizer_destroy(textRecognizer)

	img, err := packages.LoadImageMM(imgPath)
	if err != nil {
		fmt.Println(err)
	}

	status1 = C._mmdeploy_text_detector_apply(
		textDetector,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		&bboxes,
		&resultCountBbox)

	status2 = C._mmdeploy_text_recognizer_apply_bbox(
		textRecognizer,
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Max.Y), C.int(img.Rect.Max.X),
		bboxes,
		resultCountBbox,
		&texts)

	if int(status1) == 0 {
		fmt.Println("It has been sucessfully applied the model.")
	}
	if int(status2) == 0 {
		fmt.Println("It has been sucessfully applied the model.")
	}
	defer C.mmdeploy_text_recognizer_release_result((*C.struct_mmdeploy_text_recognition_t)(texts), *resultCountBbox)
	defer C.mmdeploy_text_detector_release_result((*C.struct_mmdeploy_text_detection_t)(bboxes), resultCountBbox, C.int(1))

	resultBboxes := (*[1 << 28]C.mmdeploy_text_detection_t)(unsafe.Pointer(bboxes))[:*resultCountBbox]
	resultTexts := (*[1 << 28]C.mmdeploy_text_recognition_t)(unsafe.Pointer(texts))[:*resultCountBbox]

	context := gg.NewContextForImage(img)

	for i := 0; i < len(resultBboxes); i++ {
		bboxXY := resultBboxes[i].bbox
		singleText := &resultTexts[i].text
		fmt.Printf("text: %v bbox[%v] \n", **singleText, i)
		points := CovertBboxXY2Points(bboxXY)
		*context = packages.DrawBBoxMM(*context, points)
		context.SetRGBA(0, 1, 0, 1)
		context.SetLineWidth(3)
		context.Stroke()
	}
	context.SavePNG(outputPath)
}
