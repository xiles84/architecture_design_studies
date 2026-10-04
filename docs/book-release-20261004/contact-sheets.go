// Layout proof: assemble numbered page thumbnails using only Go's image packages.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var digits = []string{"111101101101111","010110010010111","111001111100111","111001111001111","101101111001001","111100111001111","111100111101111","111001001001001","111101111101111","111101111001111"}

func main() {
	paths, err := filepath.Glob(filepath.Join(os.Args[1], "page-*.png")); must(err);sort.Strings(paths)
	for offset:=0;offset<len(paths);offset+=15 {
		canvas:=image.NewRGBA(image.Rect(0,0,1400,3430))
		draw.Draw(canvas,canvas.Bounds(),image.NewUniform(color.RGBA{225,229,235,255}),image.Point{},draw.Src)
		for i:=offset;i<len(paths)&&i<offset+15;i++ {
			f,err:=os.Open(paths[i]);must(err);page,err:=png.Decode(f);f.Close();must(err)
			x,y:=10+((i-offset)%3)*465,25+((i-offset)/3)*685
			draw.Draw(canvas,image.Rect(x,y,x+page.Bounds().Dx(),y+page.Bounds().Dy()),page,page.Bounds().Min,draw.Src)
			for j,d:=range strconv.Itoa(i+1) {
				for bit,v:=range digits[int(d-'0')] { if v=='1' {r:=image.Rect(x+j*20+(bit%3)*5,y-21+(bit/3)*3,x+j*20+(bit%3)*5+4,y-21+(bit/3)*3+3);draw.Draw(canvas,r,image.NewUniform(color.Black),image.Point{},draw.Src)} }
			}
		}
		out,err:=os.Create(filepath.Join(os.Args[2],fmt.Sprintf("contact-%02d.png",offset/15+1)));must(err);must(png.Encode(out,canvas));must(out.Close())
	}
}
func must(err error) {if err!=nil {panic(err)}}
