package captcha

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/big"
	"strconv"
)

// GenerateCode 生成指定长度的数字验证码
func GenerateCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += strconv.Itoa(int(num.Int64()))
	}
	return code
}

// GenerateImage 生成验证码图片
func GenerateImage(code string) (string, error) {
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置背景色为白色
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{240, 240, 240, 255}}, image.Point{}, draw.Src)

	// 添加噪点
	addNoise(img, 50)

	// 绘制验证码文字
	charWidth := width / len(code)
	for i, char := range code {
		x := i*charWidth + 10
		y := 15
		drawChar(img, string(char), x, y)
	}

	// 添加干扰线
	addLines(img, 3)

	// 将图片编码为PNG格式
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	// 转换为base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64Str), nil
}

// addNoise 添加噪点
func addNoise(img *image.RGBA, count int) {
	bounds := img.Bounds()
	for i := 0; i < count; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))
		r, _ := rand.Int(rand.Reader, big.NewInt(256))
		g, _ := rand.Int(rand.Reader, big.NewInt(256))
		b, _ := rand.Int(rand.Reader, big.NewInt(256))
		img.Set(int(x.Int64()), int(y.Int64()), color.RGBA{uint8(r.Int64()), uint8(g.Int64()), uint8(b.Int64()), 255})
	}
}

// addLines 添加干扰线
func addLines(img *image.RGBA, count int) {
	bounds := img.Bounds()
	for i := 0; i < count; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))

		drawLine(img, int(x1.Int64()), int(y1.Int64()), int(x2.Int64()), int(y2.Int64()))
	}
}

// drawLine 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int) {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, color.RGBA{100, 100, 100, 255})
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// drawChar 绘制字符（简化的像素字体）
func drawChar(img *image.RGBA, char string, x, y int) {
	// 简化的数字字体模式（7x9像素）
	patterns := map[string][][]int{
		"0": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
		"1": {
			{0, 0, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{1, 1, 1, 1, 1},
		},
		"2": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 0, 0, 0},
			{1, 0, 0, 0, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 1, 1, 1},
		},
		"3": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
		"4": {
			{0, 0, 0, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 1, 0, 1, 0},
			{1, 0, 0, 1, 0},
			{1, 1, 1, 1, 1},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 1, 0},
		},
		"5": {
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0},
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 0},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
		"6": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
		"7": {
			{1, 1, 1, 1, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 0, 0, 0},
		},
		"8": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
		"9": {
			{0, 1, 1, 1, 0},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 1},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0},
		},
	}

	pattern, exists := patterns[char]
	if !exists {
		return
	}

	// 添加随机颜色
	r, _ := rand.Int(rand.Reader, big.NewInt(128))
	g, _ := rand.Int(rand.Reader, big.NewInt(128))
	b, _ := rand.Int(rand.Reader, big.NewInt(128))
	charColor := color.RGBA{uint8(r.Int64()), uint8(g.Int64()), uint8(b.Int64()), 255}

	// 绘制字符
	for row, line := range pattern {
		for col, pixel := range line {
			if pixel == 1 {
				px := x + col*2
				py := y + row*2
				// 绘制2x2像素块使字符更清晰
				for dx := 0; dx < 2; dx++ {
					for dy := 0; dy < 2; dy++ {
						if px+dx < img.Bounds().Max.X && py+dy < img.Bounds().Max.Y {
							img.Set(px+dx, py+dy, charColor)
						}
					}
				}
			}
		}
	}
}
