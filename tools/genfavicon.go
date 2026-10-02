// 生成 web/favicon.ico：蓝色渐变对话气泡 + 白色省略号，与应用主题色一致。
// 输出多尺寸 ICO（256/128/64/48/32/16），每个尺寸按 4x 超采样绘制再降采样，边缘平滑。
// 用法：go run ./tools/genfavicon.go
package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

var topColor = [3]float64{111, 178, 255}   // #6fb2ff
var bottomColor = [3]float64{52, 118, 242} // #3476f2

func main() {
	sizes := []int{256, 128, 64, 48, 32, 16}
	type entry struct {
		size int
		data []byte
	}
	var entries []entry
	for _, s := range sizes {
		img := downsample(render(s*4), 4)
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			log.Fatal(err)
		}
		entries = append(entries, entry{s, buf.Bytes()})
	}

	var ico bytes.Buffer
	_ = binary.Write(&ico, binary.LittleEndian, uint16(0))
	_ = binary.Write(&ico, binary.LittleEndian, uint16(1))
	_ = binary.Write(&ico, binary.LittleEndian, uint16(len(entries)))
	offset := 6 + 16*len(entries)
	for _, e := range entries {
		b := byte(0) // 0 表示 256
		if e.size < 256 {
			b = byte(e.size)
		}
		ico.WriteByte(b)
		ico.WriteByte(b)
		ico.WriteByte(0) // 调色板色数
		ico.WriteByte(0) // 保留
		_ = binary.Write(&ico, binary.LittleEndian, uint16(1))
		_ = binary.Write(&ico, binary.LittleEndian, uint16(32))
		_ = binary.Write(&ico, binary.LittleEndian, uint32(len(e.data)))
		_ = binary.Write(&ico, binary.LittleEndian, uint32(offset))
		offset += len(e.data)
	}
	for _, e := range entries {
		ico.Write(e.data)
	}
	if err := os.WriteFile("web/favicon.ico", ico.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	// 输出预览图供人工检查
	pv := downsample(render(1024), 4)
	pf, err := os.Create(".favicon-preview.png")
	if err != nil {
		log.Fatal(err)
	}
	if err = png.Encode(pf, pv); err != nil {
		log.Fatal(err)
	}
	_ = pf.Close()
	log.Println("web/favicon.ico 生成完成，预览: .favicon-preview.png")
}

// render 以尺寸 S（正方形）绘制图标
func render(S int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, S, S))
	p := float64(S) * 0.09     // 内边距
	tailH := float64(S) * 0.15 // 气泡尾巴高度
	x0, y0 := p, p
	x1, y1 := float64(S)-p, float64(S)-p-tailH
	w := x1 - x0
	r := w * 0.24

	for py := 0; py < S; py++ {
		fy := float64(py) + 0.5
		for px := 0; px < S; px++ {
			fx := float64(px) + 0.5

			// 气泡主体（圆角矩形），垂直渐变
			bubble := inRoundRect(fx, fy, x0, y0, x1, y1, r)
			// 尾巴（左下角小三角）
			tail := inTriangle(fx, fy,
				x0+w*0.12, y1-1,
				x0+w*0.40, y1-1,
				x0+w*0.12, y1+tailH)

			if !bubble && !tail {
				continue
			}

			t := (fy - y0) / (y1 - y0)
			if t > 1 {
				t = 1
			}
			c := lerp(topColor, bottomColor, t)
			img.Set(px, py, color.RGBA{uint8(c[0]), uint8(c[1]), uint8(c[2]), 255})
		}
	}

	// 三个白色省略号圆点
	cyb := (y0 + y1) / 2
	dotR := float64(S) * 0.055
	for _, fx := range []float64{x0 + w*0.27, x0 + w*0.5, x0 + w*0.73} {
		for py := 0; py < S; py++ {
			for px := 0; px < S; px++ {
				dx := float64(px) + 0.5 - fx
				dy := float64(py) + 0.5 - cyb
				if dx*dx+dy*dy <= dotR*dotR {
					img.Set(px, py, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
	return img
}

func inRoundRect(x, y, x0, y0, x1, y1, r float64) bool {
	if x < x0 || x > x1 || y < y0 || y > y1 {
		return false
	}
	corners := [][2]float64{{x0 + r, y0 + r}, {x1 - r, y0 + r}, {x0 + r, y1 - r}, {x1 - r, y1 - r}}
	for _, c := range corners {
		cx, cy := c[0], c[1]
		outsideX := (cx < (x0+x1)/2 && x < cx) || (cx > (x0+x1)/2 && x > cx)
		outsideY := (cy < (y0+y1)/2 && y < cy) || (cy > (y0+y1)/2 && y > cy)
		if outsideX && outsideY {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy > r*r {
				return false
			}
		}
	}
	return true
}

func inTriangle(px, py, ax, ay, bx, by, cx, cy float64) bool {
	d1 := sign(px, py, ax, ay, bx, by)
	d2 := sign(px, py, bx, by, cx, cy)
	d3 := sign(px, py, cx, cy, ax, ay)
	hasNeg := d1 < 0 || d2 < 0 || d3 < 0
	hasPos := d1 > 0 || d2 > 0 || d3 > 0
	return !(hasNeg && hasPos)
}

func sign(px, py, ax, ay, bx, by float64) float64 {
	return (px-bx)*(ay-by) - (ax-bx)*(py-by)
}

func lerp(a, b [3]float64, t float64) [3]float64 {
	return [3]float64{
		math.Round(a[0] + (b[0]-a[0])*t),
		math.Round(a[1] + (b[1]-a[1])*t),
		math.Round(a[2] + (b[2]-a[2])*t),
	}
}

// downsample 按 n*n 块平均降采样（预乘 alpha，避免透明边缘出现黑边）
func downsample(src *image.RGBA, n int) *image.RGBA {
	S := src.Bounds().Dx()
	D := S / n
	dst := image.NewRGBA(image.Rect(0, 0, D, D))
	for dy := 0; dy < D; dy++ {
		for dx := 0; dx < D; dx++ {
			var r, g, b, a float64
			for j := 0; j < n; j++ {
				for i := 0; i < n; i++ {
					c := src.RGBAAt(dx*n+i, dy*n+j)
					fr := float64(c.R) / 255
					fg := float64(c.G) / 255
					fb := float64(c.B) / 255
					fa := float64(c.A) / 255
					r += fr * fa
					g += fg * fa
					b += fb * fa
					a += fa
				}
			}
			cnt := float64(n * n)
			r /= cnt
			g /= cnt
			b /= cnt
			a /= cnt
			var out color.RGBA
			if a > 0 {
				out = color.RGBA{
					uint8(math.Round(r / a * 255)),
					uint8(math.Round(g / a * 255)),
					uint8(math.Round(b / a * 255)),
					uint8(math.Round(a * 255)),
				}
			}
			dst.SetRGBA(dx, dy, out)
		}
	}
	return dst
}
