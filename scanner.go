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
	"unsafe"
)

type SymbolType int

const (
	All     SymbolType = 0
	QR      SymbolType = C.ZBAR_QRCODE
	EAN13   SymbolType = C.ZBAR_EAN13
	EAN8    SymbolType = C.ZBAR_EAN8
	UPCA    SymbolType = C.ZBAR_UPCA
	UPCE    SymbolType = C.ZBAR_UPCE
	ISBN10  SymbolType = C.ZBAR_ISBN10
	ISBN13  SymbolType = C.ZBAR_ISBN13
	Code128 SymbolType = C.ZBAR_CODE128
	Code39  SymbolType = C.ZBAR_CODE39
	PDF417  SymbolType = C.ZBAR_PDF417
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

// Close 메모리 해제
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

// Scan 그레이스케일 데이터 스캔
func (s *Scanner) Scan(data []byte, width int, height int) ([]BarcodeData, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	if s.scannerPtr == nil || s.imagePtr == nil {
		return nil, errors.New("scanner is nil")
	}
	C.zbar_image_set_size(s.imagePtr, C.uint(width), C.uint(height))
	cData := unsafe.Pointer(&data[0])
	C.zbar_image_set_data(s.imagePtr, cData, C.ulong(len(data)), nil)

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
