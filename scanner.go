package zgobar

// 1.0

/*
#cgo LDFLAGS: -lzbar
#include <zbar.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"image"
	"unsafe"
)

type SymbolType int

const (
	All        SymbolType = 0
	QR         SymbolType = C.ZBAR_QRCODE
	EAN13      SymbolType = C.ZBAR_EAN13
	EAN8       SymbolType = C.ZBAR_EAN8
	UPCA       SymbolType = C.ZBAR_UPCA
	UPCE       SymbolType = C.ZBAR_UPCE
	ISBN10     SymbolType = C.ZBAR_ISBN10
	ISBN13     SymbolType = C.ZBAR_ISBN13
	Code128    SymbolType = C.ZBAR_CODE128
	Code39     SymbolType = C.ZBAR_CODE39
	PDF417     SymbolType = C.ZBAR_PDF417
	I25        SymbolType = C.ZBAR_I25
	Code93     SymbolType = C.ZBAR_CODE93
	Codabar    SymbolType = C.ZBAR_CODABAR
	DataBar    SymbolType = C.ZBAR_DATABAR
	DataBarExp SymbolType = C.ZBAR_DATABAR_EXP
	Aztec      SymbolType = C.ZBAR_AZTEC
)

type BarcodeData struct {
	Type string
	Data string
}

type Scanner struct {
	scannerPtr *C.zbar_image_scanner_t
	imagePtr   *C.zbar_image_t
}

func NewScanner(types ...SymbolType) *Scanner {
	s := C.zbar_image_scanner_create()
	if len(types) == 0 {
		C.zbar_image_scanner_set_config(s, 0, C.ZBAR_CFG_ENABLE, 1)
	} else {
		C.zbar_image_scanner_set_config(s, 0, C.ZBAR_CFG_ENABLE, 0)
		for _, t := range types {
			C.zbar_image_scanner_set_config(s, C.zbar_symbol_type_t(t), C.ZBAR_CFG_ENABLE, 1)
		}
	}
	img := C.zbar_image_create()
	// 8비트 흑백 이미지 설정
	C.zbar_image_set_format(img, C.fourcc('Y', '8', '0', '0'))
	return &Scanner{
		scannerPtr: s,
		imagePtr:   img,
	}
}

func (s *Scanner) Close() {
	if s.imagePtr != nil {
		C.zbar_image_destroy(s.imagePtr)
		s.imagePtr = nil
	}
	if s.scannerPtr != nil {
		C.zbar_image_scanner_destroy(s.scannerPtr)
		s.scannerPtr = nil
	}
}

func (s *Scanner) Scan(img image.Image) ([]BarcodeData, error) {
	if img == nil {
		return nil, errors.New("image is nil")
	}
	if s.scannerPtr == nil || s.imagePtr == nil {
		return nil, errors.New("scanner is closed")
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width == 0 || height == 0 {
		return nil, errors.New("image size is 0")
	}

	C.zbar_image_set_size(s.imagePtr, C.uint(width), C.uint(height))

	var cData unsafe.Pointer
	var rawBytes []byte // GC 방지용

	switch v := img.(type) {
	case *image.Gray:
		cData = unsafe.Pointer(&v.Pix[0])
	default:
		// 컬러면 흑백으로 변환
		rawBytes = toGrayBytes(img, width, height)
		cData = unsafe.Pointer(&rawBytes[0])
	}

	C.zbar_image_set_data(s.imagePtr, cData, C.ulong(width*height), nil)
	n := C.zbar_scan_image(s.scannerPtr, s.imagePtr)
	if n < 0 {
		return nil, errors.New("scan error")
	}
	if n == 0 {
		return nil, nil
	}

	var results []BarcodeData
	symbol := C.zbar_image_first_symbol(s.imagePtr)
	for symbol != nil {
		res := BarcodeData{
			Type: C.GoString(C.zbar_get_symbol_name(C.zbar_symbol_get_type(symbol))),
			Data: C.GoString(C.zbar_symbol_get_data(symbol)),
		}
		results = append(results, res)
		symbol = C.zbar_symbol_next(symbol)
	}
	return results, nil
}

func toGrayBytes(img image.Image, w, h int) []byte {
	gray := make([]byte, w*h)
	if rgbaImg, ok := img.(*image.RGBA); ok {
		for i, idx := 0, 0; i < len(rgbaImg.Pix); i += 4 {
			r := uint32(rgbaImg.Pix[i])
			g := uint32(rgbaImg.Pix[i+1])
			b := uint32(rgbaImg.Pix[i+2])
			gray[idx] = byte((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			idx++
		}
		return gray
	}
	if ycbcrImg, ok := img.(*image.YCbCr); ok {
		if len(ycbcrImg.Y) == w*h {
			copy(gray, ycbcrImg.Y)
			return gray
		}
	}
	// 기타 이미지
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			gray[y*w+x] = byte((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
		}
	}
	return gray
}
