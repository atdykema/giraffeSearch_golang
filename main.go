package main

import (
	"fmt"
	"log"
	"os"
	//"os/exec"
	//"strconv"
	"strings"
	"sync"
	"bufio"
)

var MAX_DEPTH int
var output_message_err []string
//var output_message_files []string
var cGUI chan string
var wg *sync.WaitGroup

func init(){
	//init config
	MAX_DEPTH = 3
	output_message_err = make([]string, 0)
	//output_message_files = make([]string, 0)
	cGUI = make(chan string)
	wg = &sync.WaitGroup{}
}

func main(){
	
	keyword := getKeyword()
	pwd := "/" //allow manual input
	count := 0
	depth := 0 //allow manual input

	wg.Add(1)
	go startFileSearch(keyword, pwd, count, depth)

	wg.Add(1)
	go callCLIGUI()
	
	wg.Wait()

	if len(output_message_err) == 0{
		fmt.Println("No Errors")
	}else{
		file, err := os.OpenFile("errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		err = os.Truncate("errors.txt", 0)

		if err != nil {
			log.Fatalf("failed truncating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range output_message_err {
			_, _ = datawriter.WriteString(data + "\n")
		}
	 
		datawriter.Flush()
		file.Close()

		/*
		for _, m := range output_message_err{
			fmt.Println(m)
		}
		*/
	}
}

func startFileSearch(keyword string, pwd string, count int, depth int){
	//wg.Add(1)
	searchFile(keyword, pwd, count, depth)
	//wg.Done()
	//end := make([]string, 0)
	//end = append(end, "END")
	cGUI <- "END"

	wg.Done()
}

func callCLIGUI(){
	payload := <- cGUI
	//for payload[0] != "END"{
	for payload != "END"{
		payload = <- cGUI
		fmt.Println(payload)
		/*
		for _, ele := range payload{
			fmt.Println(ele)
		}
		*/

	}

	close(cGUI)

	wg.Done()

}

/*
func callClear(){
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
*/

func getKeyword() (keyword string){
	keyword = ""
	fmt.Println("Enter keyword to search for: ")
	fmt.Scanln(&keyword)

	return keyword
}

func searchFile(keyword string, pwd string, count int, depth int){
	
	files, err := os.ReadDir(pwd)

	if err != nil {
		message := strings.Split(err.Error(), ": ")
		if message[len(message)-1] == "permission denied"{
			output_message_err = append(output_message_err, err.Error())
			return
		}else if message[len(message)-1] == "operation not permitted"{
			output_message_err = append(output_message_err, err.Error())
			return
		}else if message[len(message)-1] == "bad file descriptor"{
			output_message_err = append(output_message_err, err.Error())
			return
		}else{
			log.Fatal(err)
		}
    	
	}

	for _, f := range files {
		//fmt.Println(1)

		if f.Name()[0] == 46{ //dot
			continue
		}

		count++
		extraSlash := ""

		if pwd != "/"{
			extraSlash = "/"
		}
		//fmt.Println(2)
		fstats, err := os.Stat(pwd + extraSlash + f.Name())
		//fmt.Println(3)
		if err != nil {
			//fmt.Println(4)
			message := strings.Split(err.Error(), ": ")
			if message[len(message)-1] == "file name too long"{
				continue
			}else if message[len(message)-1] == "no such file or directory"{
				output_message_err = append(output_message_err, err.Error())
				continue
			}else if message[len(message)-1] == "permission denied"{
				output_message_err = append(output_message_err, err.Error())
				continue
			}else if message[len(message)-1] == "operation not permitted"{
				output_message_err = append(output_message_err, err.Error())
				continue
			}else if message[len(message)-1] == "bad file descriptor"{
				output_message_err = append(output_message_err, err.Error())
				continue
			}else{
				//fmt.Println(5)
				log.Fatal(err)
			}
		}else{
			//fmt.Println(6)
			curr_pwd := pwd + (extraSlash + fstats.Name())
			//fmt.Println(curr_pwd)
			if keyword == fstats.Name(){
				//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
				//output_message_files = append(output_message_files, curr_pwd)
				cGUI <- curr_pwd

			}else if fstats.IsDir(){
				//fmt.Println(7)
				if depth > MAX_DEPTH{
					continue
				}
				//fmt.Println(8)
				searchFile(keyword, curr_pwd, count, depth+1)
			}
		}
	}

}


