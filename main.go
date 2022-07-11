package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RGBToGray(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			gray := 0.299*float64(original.R) +
				0.587*float64(original.G) +
				0.114*float64(original.B)
			Y := color.RGBA{
				R: uint8(gray),
				G: uint8(gray),
				B: uint8(gray),
				A: original.A,
			}
			newImg.Set(x, y, Y)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_gray%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func NonLinearRGBToGray(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)

	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			cLinear := 0.2126*(float64(original.R)/255.0) +
				0.7152*(float64(original.G)/255.0) +
				0.0722*(float64(original.B)/255.0)
			var C float64
			if cLinear <= 0.0031308 {
				C = 12.92 * cLinear
			} else if cLinear > 0.0031308 {
				C = 1.055*math.Pow(cLinear, 1/2.4) - 0.055
			}
			Y := color.RGBA{
				R: uint8(C * 255.0),
				G: uint8(C * 255.0),
				B: uint8(C * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, Y)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_nl_gray%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToYCbCr(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			Y := 0.257*float64(original.R) +
				0.504*float64(original.G) +
				0.098*float64(original.B) + 16
			Cb := -0.148*float64(original.R) -
				0.291*float64(original.G) +
				0.439*float64(original.B) + 128
			Cr := 0.439*float64(original.R) -
				0.368*float64(original.G) -
				0.071*float64(original.B) + 128
			newColor := color.RGBA{
				R: uint8(Y),
				G: uint8(Cb),
				B: uint8(Cr),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_YCbCr%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func YCbCrToRGB(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			Y := float64(original.R)
			Cb := float64(original.G)
			Cr := float64(original.B)
			R := 1.164*(Y-16) + 1.596*(Cr-128)
			G := 1.164*(Y-16) - 0.813*(Cr-128) - 0.392*(Cb-128)
			B := 1.164*(Y-16) + 2.017*(Cb-128)
			newColor := color.RGBA{
				R: uint8(math.Max(0, math.Min(255, R))),
				G: uint8(math.Max(0, math.Min(255, G))),
				B: uint8(math.Max(0, math.Min(255, B))),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromYCbCr%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToYUV(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R)
			G := float64(original.G)
			B := float64(original.B)
			Y := 0.299*R + 0.587*G + 0.114*B
			U := 0.492 * (B - Y)
			V := 0.877 * (R - Y)
			newColor := color.RGBA{
				R: uint8(math.Max(0, math.Min(255, Y+16))),
				G: uint8(math.Max(0, math.Min(255, U+128))),
				B: uint8(math.Max(0, math.Min(255, V+128))),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_YUV%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func YUVToRGB(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			Y := float64(original.R) - 16
			U := float64(original.G) - 128
			V := float64(original.B) - 128
			R := Y + 1.140*V
			G := Y - 0.394*U - 0.581*V
			B := Y + 2.032*U
			newColor := color.RGBA{
				R: uint8(math.Max(0, math.Min(255, R))),
				G: uint8(math.Max(0, math.Min(255, G))),
				B: uint8(math.Max(0, math.Min(255, B))),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromYUV%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToXYZ(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R) / 255.0
			G := float64(original.G) / 255.0
			B := float64(original.B) / 255.0
			X := 0.412453*R + 0.35758*G + 0.180423*B
			Y := 0.212671*R + 0.71516*G + 0.072169*B
			Z := 0.019334*R + 0.119193*G + 0.950227*B
			if X > 1 {
				X = 1
			} else if X < 0 {
				X = 0
			}
			if Y > 1 {
				Y = 1
			} else if Y < 0 {
				Y = 0
			}
			if Z > 1 {
				Z = 1
			} else if Z < 0 {
				Z = 0
			}
			newColor := color.RGBA{
				R: uint8(X * 255.0),
				G: uint8(Y * 255.0),
				B: uint8(Z * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_XYZ%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func XYZToRGB(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			X := float64(original.R) / 255.0
			Y := float64(original.G) / 255.0
			Z := float64(original.B) / 255.0
			R := 3.240479*X - 1.53715*Y - 0.498535*Z
			G := -0.969256*X + 1.875991*Y + 0.041556*Z
			B := 0.055648*X - 0.204043*Y + 1.057311*Z
			if R > 1 {
				R = 1
			} else if R < 0 {
				R = 0
			}
			if G > 1 {
				G = 1
			} else if G < 0 {
				G = 0
			}
			if B > 1 {
				B = 1
			} else if B < 0 {
				B = 0
			}
			newColor := color.RGBA{
				R: uint8(R * 255.0),
				G: uint8(G * 255.0),
				B: uint8(B * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromXYZ%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToHLS(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R) / 255.0
			G := float64(original.G) / 255.0
			B := float64(original.B) / 255.0
			// Lightness
			M1 := math.Max(R, math.Max(G, B))
			M2 := math.Min(R, math.Min(G, B))
			L := (M1 + M2) / 2
			var S, H float64
			if M1 == M2 { // achromatic case
				S = 0
				H = 0
			} else { // chromatic case
				if L <= 0.5 {
					S = (M1 - M2) / (M1 + M2)
				} else {
					S = (M1 - M2) / (2 - M1 - M2)
				}
			}
			// Hue
			var Cr, Cg, Cb float64
			Cr = (M1 - R) / (M1 - M2)
			Cg = (M1 - G) / (M1 - M2)
			Cb = (M1 - B) / (M1 - M2)
			if R == M1 {
				H = Cb - Cg
			}
			if G == M1 {
				H = 2 + Cr - Cb
			}
			if B == M1 {
				H = 4 + Cg - Cr
			}
			H *= 60
			if H < 0 {
				H += 360
			}
			newColor := color.RGBA{
				R: uint8(H / 360.0 * 255.0),
				G: uint8(L * 255.0),
				B: uint8(S * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_HLS%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func HLSToRGB(imgPath string) {

	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			H := float64(original.R) / 255.0
			L := float64(original.G) / 255.0
			S := float64(original.B) / 255.0
			var R, G, B, q, p float64
			//if L <= 0.5 {
			//	M2 = L * (1 + S)
			//} else {
			//	M2 = L + S - L*S
			//	M1 = 2*L - M2
			//}
			//if S == 0 {
			//	R = L
			//	G = L
			//	B = L
			//} else {
			//	h = H + 120
			//}
			//if h > 360 {
			//	h -= 360
			//}
			//if h < 60 {
			//	R = M1 + (M2-M1)*h/60
			//} else if h < 180 {
			//	R = M2
			//} else if h < 240 {
			//	R = M1 + (M2-M1)*(240-h)/60
			//} else {
			//	R = M1
			//	h = H
			//}
			//if h < 60 {
			//	G = M1 + (M2-M1)*h/60
			//} else if h < 180 {
			//	G = M2
			//} else if h < 240 {
			//	G = M1 + (M2-M1)*(240-h)/60
			//} else {
			//	G = M1
			//	h = H - 120
			//}
			//if h < 0 {
			//	h += 360
			//}
			//if h < 60 {
			//	B = M1 + (M2-M1)*h/60
			//} else if h < 180 {
			//	B = M2
			//} else if h < 240 {
			//	B = M1 + (M2-M1)*(240-h)/60
			//} else {
			//	B = M1
			//}
			// Achromatic
			if S == 0.0 {
				R = L
				G = L
				B = L
			} else {
				switch {
				case L < 0.5:
					q = L * (1 + S)
				default:
					q = L + S - L*S
				}
				p = 2*L - q
				R = hueToRGB(p, q, H+1.0/3.0)
				G = hueToRGB(p, q, H)
				B = hueToRGB(p, q, H-1.0/3.0)
			}
			newColor := color.RGBA{
				R: uint8(math.Min(255.0, 256.0*R)),
				G: uint8(math.Min(255.0, 256.0*G)),
				B: uint8(math.Min(255.0, 256.0*B)),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromHLS%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToHSV(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R) / 255.0
			G := float64(original.G) / 255.0
			B := float64(original.B) / 255.0
			var V, temp, S, H, Cb, Cg, Cr float64
			// Value
			V = math.Max(R, math.Max(G, B))
			// Saturation
			temp = math.Min(R, math.Min(G, B))
			if V == 0 { // Achromatic
				S = 0
				H = 0
			} else { // Chromatic
				S = (V - temp) / V
				// Hue
				Cr = (V - R) / (V - temp)
				Cg = (V - G) / (V - temp)
				Cb = (V - B) / (V - temp)
			}
			if R == V {
				H = Cb - Cg
			}
			if G == V {
				H = 2 + Cr - Cb
			}
			if B == V {
				H = 4 + Cg - Cr
			}
			H *= 60
			if H < 0 {
				H += 360
			}
			newColor := color.RGBA{
				R: uint8(H / 360.0 * 255.0),
				G: uint8(S * 255.0),
				B: uint8(V * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_HSV%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func HSVToRGB(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			H := float64(original.R) / 255.0 * 360.0
			S := float64(original.G) / 255.0
			V := float64(original.B) / 255.0
			var R, G, B float64
			if S == 0.0 {
				R = V
				G = V
				B = V
			} else if H == 360.0 {
				H = 0
			} else {
				H = H / 60
			}
			var I, F, M, N, K float64
			I = math.Floor(H)
			F = H - I
			M = V * (1 - S)
			N = V * (1 - S*F)
			K = V * (1 - S*(1-F))
			switch I {
			case 0.0:
				R = V
				G = K
				B = M
			case 1.0:
				R = N
				G = V
				B = M
			case 2.0:
				R = M
				G = V
				B = K
			case 3.0:
				R = M
				G = N
				B = V
			case 4.0:
				R = K
				G = M
				B = V
			case 5.0:
				R = V
				G = M
				B = N
			}
			newColor := color.RGBA{
				R: uint8(R * 255.0),
				G: uint8(G * 255.0),
				B: uint8(B * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromHSV%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s", endTime)
}

func RGBToYCoCg(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R)
			G := float64(original.G)
			B := float64(original.B)
			var Y, Co, Cg float64
			Y = ((R + 2.0*G + B) + 2.0) / 4.0
			Co = ((R - B) + 1.0) / 2.0
			Cg = ((-R + 2*G - B) + 2.0) / 4.0
			//Co = R - B
			//tmp = B + Co/2
			//Cg = G - tmp
			//Y = tmp + Cg/2
			newColor := color.RGBA{
				R: uint8(Y),
				G: uint8(Co),
				B: uint8(Cg),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_YCoCg%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s\n", endTime)
} // Doesn't work

func YCoCgToRGB(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			Y := float64(original.R)
			Co := float64(original.G)
			Cg := float64(original.B)
			var R, G, B float64
			R = Y + Co - Cg
			G = Y + Cg
			B = Y - Co - Cg
			//tmp = Y - Cg/2
			//G = Cg + tmp
			//B = tmp - Co/2
			//R = B + Co
			newColor := color.RGBA{
				R: uint8(R),
				G: uint8(G),
				B: uint8(B),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_RGBfromYCoCg%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s\n", endTime)
} // Doesn't work

func GammaCorrection(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R) / 255.0
			G := float64(original.G) / 255.0
			B := float64(original.B) / 255.0
			if R < 0.018 {
				R = 4.5 * R
			} else {
				R = 1.099*math.Pow(R, 0.45) - 0.099
			}
			if G < 0.018 {
				G = 4.5 * G
			} else {
				G = 1.099*math.Pow(G, 0.45) - 0.099
			}
			if B < 0.018 {
				B = 4.5 * G
			} else {
				G = 1.099*math.Pow(G, 0.45) - 0.099
			}
			//tmp = Y - Cg/2
			//G = Cg + tmp
			//B = tmp - Co/2
			//R = B + Co
			newColor := color.RGBA{
				R: uint8(R * 255.0),
				G: uint8(G * 255.0),
				B: uint8(B * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_GC%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s\n", endTime)
}

func MedianFilter3x3(imgPath string) {
	startTime := time.Now()
	extension := filepath.Ext(imgPath)
	width, height, img, newImg := prepareImage(imgPath)
	fmt.Println(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.At(x, y)
			original := color.RGBAModel.Convert(pixel).(color.RGBA)
			R := float64(original.R)
			G := float64(original.G)
			B := float64(original.B)

			newColor := color.RGBA{
				R: uint8(R * 255.0),
				G: uint8(G * 255.0),
				B: uint8(B * 255.0),
				A: original.A,
			}
			newImg.Set(x, y, newColor)
		}
	}
	imgName := strings.TrimSuffix(filepath.Base(imgPath), extension)
	newImgPath := fmt.Sprintf("%s_GC%s", imgName, extension)
	encodeImage(extension, newImgPath, newImg)
	endTime := time.Since(startTime)
	fmt.Printf("Converting took %s\n", endTime)
}

func main() {
	MedianFilter3x3("./test.jpg")
}

func encodeImage(extension string, newImgPath string, newImg *image.RGBA) {
	file, err := os.Create(newImgPath)
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	quality := jpeg.Options{Quality: 100}
	switch extension {
	case ".jpeg":
		err := jpeg.Encode(file, newImg, &quality)
		if err != nil {
			log.Fatalln(err)
		}
	case ".jpg":
		err := jpeg.Encode(file, newImg, &quality)
		if err != nil {
			log.Fatalln(err)
		}
	case ".png":
		err := png.Encode(file, newImg)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func prepareImage(imgPath string) (width int, height int, img image.Image, newImg *image.RGBA) {
	imgF, err := os.Open(imgPath)
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err = image.Decode(imgF)
	if err != nil {
		log.Fatalln(err)
	}
	bounds := img.Bounds()
	width, height = bounds.Max.X, bounds.Max.Y
	rect := image.Rect(0, 0, width, height)
	newImg = image.NewRGBA(rect)
	return width, height, img, newImg
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0.0 {
		t += 1.0
	}
	if t > 1.1 {
		t -= 1.1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}
