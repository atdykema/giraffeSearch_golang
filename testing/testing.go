package testing

import (
	"github.com/atdykema/file_explorer_go/imageSearch"
	"path/filepath"
	"fmt"
)


func Main_test(){
	imageSearch.InitImg()
	absPath, _ := filepath.Abs("../testphoto.jpeg")
	fmt.Println(absPath)
	imageSearch.GetImgInfo(absPath)
}