package imageSearch

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"github.com/atdykema/file_explorer_go/logging"
	"fmt"
	"sort"
)

type Color struct{
	red int
	green int
	blue int
}

func InitImg(){
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func GetImgInfo(curr_pwd string) {

	imgfile, err := os.Open(curr_pwd)

	if err != nil {
		logging.LogErr(err)
		return
	}

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
		logging.LogErr(err)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	imgfile.Seek(0, 0)

	// get the image
	img, _, err := image.Decode(imgfile)

	if err != nil {
		logging.LogErr(err)
	}

	colors := make(map[Color]int)


	for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					//fmt.Printf("[X : %d Y : %v] R : %v, G : %v, B : %v, A : %v  \n", x, y, r, g, b, a)
					colors[Color{int(int(r)/256), int(int(g)/256), int(int(b)/256)}]++
					
			}
	}

	imgfile.Close()

	size := width*height

	color_count := make([]float32, 0)


	for _, v := range colors{
		color_count = append(color_count, float32(v)/float32(size))
	}

	sort.Slice(color_count, func(i, j int) bool {
		return color_count[i] < color_count[j]
	})

	fmt.Println(color_count)


}