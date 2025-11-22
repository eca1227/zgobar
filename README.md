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
### 함수
```go
// 객체생성
NewScanner(types ...SymbolType) *Scanner

// 객체닫기
(s *Scanner) Close()

// 이미지 스캔
(s *Scanner) Scan(data []byte, width int, height int) ([]BarcodeData, error)
```

### 예제

```go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/eca1227/zgobar"
)

func main() {
	scanner := zgobar.NewScanner(zgobar.QR)
	defer scanner.Close()

	file, _ := os.Open("barcode.png")
	defer file.Close()
	img, _, _ := image.Decode(file)

	gray, w, h := toGray(img)
	
	results, err := scanner.Scan(gray, w, h)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		fmt.Println("바코드를 찾을 수 없습니다.")
	} else {
		for _, res := range results {
			fmt.Printf("[%s] %s\n", res.Type, res.Data)
		}
	}
}

// 이미지를 그레이스케일(Y800) []byte로 변환하는 헬퍼 함수
func toGray(img image.Image) ([]byte, int, int) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	gray := make([]byte, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x+b.Min.X, y+b.Min.Y)
			r, g, blue, _ := c.RGBA()
			// RGB to Grayscale 공식
			gray[y*w+x] = byte((19595*r + 38470*g + 7471*blue + 1<<15) >> 24)
		}
	}
	return gray, w, h
}
```

## License
MIT License (Linked against LGPL-2.1 ZBar)