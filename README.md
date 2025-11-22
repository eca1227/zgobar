# zgobar

**zgobar**는 Go 언어용 ZBar 래퍼. **모든 1D/2D 바코드를 자동 인식**

## 필수
이 라이브러리는 **ZBar C 라이브러리**가 시스템에 설치되어 있어야 작동합니다.

- **Ubuntu / Debian / Raspberry Pi**: `sudo apt-get install libzbar-dev`
- **macOS**: `brew install zbar`
- **CentOS / RHEL**: `sudo yum install zbar-devel`
- **Windows (MSYS2)**: `pacman -S mingw-w64-x86_64-zbar`

## 설치

```bash
go get "github.com/eca1227/zgobar"
```

## 사용법
### 상수 (SymbolType)
```go
// NewScanner()에 인자로 전달

//1D 바코드
zgobar.EAN8 - EAN-8 (유럽 상품코드)
zgobar.EAN13 - EAN-13 (유럽 상품코드)
zgobar.UPCE - UPC-E (미국 상품코드)
zgobar.UPCA - UPC-A (미국 상품코드)
zgobar.ISBN10 - ISBN-10 (책 코드)
zgobar.ISBN13 - ISBN-13 (책 코드)
zgobar.I25 - Interleaved 2 of 5
zgobar.CODE39 - Code 39
zgobar.CODE93 - Code 93
zgobar.CODE128 - Code 128
zgobar.CODABAR - Codabar

//2D 바코드
zgobar.QR - QR Code
zgobar.DATABAR - GS1 DataBar
zgobar.DATABAR_EXP - GS1 DataBar Expanded
zgobar.PDF417 - PDF417
zgobar.AZTEC - Aztec
```
### 함수
```go
// 객체생성 (주의: 고루틴 당 하나의 객체만 사용)
NewScanner(types ...SymbolType) *Scanner
// 예시
scanner := zgobar.NewScanner(zgobar.QR) // 단일 타입 스캔
scanner := zgobar.NewScanner(zgobar.QR, zgobar.CODE128, zgobar.EAN13) // 여러 타입 스캔
scanner := zgobar.NewScanner()// 모든 타입 스캔


// 객체닫기 (객체 생성 후 반드시 defer로 호출)
(s *Scanner) Close()
// 예시
defer scanner.Close()


// 이미지 스캔
// img - image.Image 인터페이스 (Gray, RGBA, YCbCr 등 모든 이미지 타입 지원)
// 반환: 인식된 바코드 데이터 슬라이스, 에러
(s *Scanner) Scan(img image.Image) ([]BarcodeData, error)

type BarcodeData struct {
Type string  // 바코드 타입명 (예: "QR-CODE", "CODE128")
Data string  // 인식된 바코드 데이터 (문자열)
}
```

### 예제

```go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"github.com/eca1227/zgobar"
)

func main() {
	// 여러 바코드 타입 스캔
	scanner := zgobar.NewScanner(zgobar.QR, zgobar.Code128, zgobar.EAN13)
	defer scanner.Close()

	// 이미지 파일 열기
	file, _ := os.Open("barcode.png")
	defer file.Close()
	img, _, _ := image.Decode(file)

	results, err := scanner.Scan(img)
	if err != nil {
		log.Fatalf("스캔 실패: %v", err)
	} else if results == nil {
		fmt.Println("바코드를 찾을 수 없습니다.")
	} else {
		for _, barcode := range results {
			fmt.Printf("[%s] %s\n", barcode.Type, barcode.Data)
		}
	}
}
```

## License
MIT License (Linked against LGPL-2.1 ZBar)