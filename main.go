package main

import (
	"fmt"
	"os"
	"log"

)

func main(){
	keyword := getKeyword()
	search(keyword)
}

func getKeyword() (keyword string){
	keyword = ""
	fmt.Println("Enter keyword to search for: ")
	fmt.Scanln(&keyword)
	
	return keyword
}

func search(keyword string){

	files, err := os.ReadDir(keyword)

	if err != nil {
    	log.Fatal(err)
	}

	for _, f := range files {
    	fmt.Println(f.Name())
	}
}


