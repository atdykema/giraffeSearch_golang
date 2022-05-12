package imageSearch

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"github.com/atdykema/file_explorer_go/logging"
	"fmt"
)

func InitImg(){
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func GetImgInfo(curr_pwd string){

	imgfile, err := os.Open(curr_pwd)

	if err != nil {
		logging.LogErr(err)
		return
	}

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)




	imgfile.Close()


}