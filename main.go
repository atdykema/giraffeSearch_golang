package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main(){
	keyword := getKeyword()
	pwd := "/" //allow manual input
	count := 0
	depth := 0 //allow manual input
	output_message := make([]string, 0)
	searchFile(keyword, pwd, count, depth, output_message)
	for _, m := range output_message{
		fmt.Println(m)
	}
}

func getKeyword() (keyword string){
	keyword = ""
	fmt.Println("Enter keyword to search for: ")
	fmt.Scanln(&keyword)

	return keyword
}

func searchFile(keyword string, pwd string, count int, depth int, output_message []string){
	
	files, err := os.ReadDir(pwd)

	if err != nil {
    	log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(1)

		if f.Name()[0] == 46{ //dot
			continue
		}

		count++
		extraSlash := ""

		if pwd != "/"{
			extraSlash = "/"
		}
		fmt.Println(2)
		fstats, err := os.Stat(pwd + extraSlash + f.Name())
		fmt.Println(3)
		if err != nil {
			fmt.Println(4)
			message := strings.Split(err.Error(), ": ")
			fmt.Println("peanutbutter", message)
			if message[len(message)-1] == "file name too long"{
				continue
			}else if message[len(message)-1] == "no such file or directory"{
				output_message = append(output_message, err.Error())
				continue
			}else if message[len(message)-1] == "permission denied"{
				fmt.Println("bruh")
				output_message = append(output_message, err.Error())
				continue
			}else{
				fmt.Println(5)
				log.Fatal(err)
			}
		}else{
			fmt.Println(6)
			curr_pwd := pwd + (extraSlash + fstats.Name())
			fmt.Println(curr_pwd)
			if keyword == fstats.Name(){
				fmt.Println(fstats.Name())
				fmt.Println(pwd)
				fmt.Println("Files searched: ", count)
				break

			}else if fstats.IsDir(){
				fmt.Println(6)
				if depth > 20{
					continue
				}
				fmt.Println(7)
				searchFile(keyword, curr_pwd, count, depth+1, output_message)
			}
			fmt.Println(fstats.Name(), "skipped")
		}
	}
}


