package main

import (
	"fmt"
	"log"
	"os"
	//"os/exec"
	//"strconv"
	"bufio"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type Configuration struct {
	CommandLineConfig bool
	MAX_DEPTH int
	SearchType int
	ExactMatch bool
	PrintConfigOnSearch bool
}

//!!!!!!!!!!!!!!!!!!!!

//TODO: allow search for file type, last modification


var output_message_err []string
var output_message_files []string
var cGUI chan string
var wg_main *sync.WaitGroup
var wg_deep *sync.WaitGroup
var stack []string
var newStack []string
var shallowDepth int
var configuration Configuration
//var mutex_search sync.Mutex

func init_config() Configuration{
	file, _ := os.Open("./config/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		log.Fatal(err.Error())
	}
	return configuration
}

func init(){
	//init config
	configuration = init_config()
	output_message_err = make([]string, 0)
	output_message_files = make([]string, 0)
	cGUI = make(chan string)
	wg_main = &sync.WaitGroup{}
	wg_deep = &sync.WaitGroup{}
	stack = []string{}
	newStack = []string{}
	shallowDepth = 0
}

func main(){
	
	if configuration.PrintConfigOnSearch{
		fmt.Println("Config:")
		fmt.Println("CommandLineConfig: ", configuration.CommandLineConfig)
		fmt.Println("MAX_DEPTH: ", configuration.MAX_DEPTH)
		fmt.Println("SearchType: ", configuration.SearchType)
		fmt.Println("ExactMatch: ", configuration.ExactMatch)
		fmt.Println("PrintConfigOnSearch: ", configuration.PrintConfigOnSearch)
	}

	keyword := getKeyword()

	initTime := time.Now()


	pwd := "/" //allow manual input
	count := 0
	depth := 0 //allow manual input

	fmt.Println("---")

	wg_main.Add(1)
	go startFileSearch(keyword, pwd, count, depth)

	wg_main.Add(1)
	go callCLIGUI()
	
	wg_main.Wait()

	elapsedTime := time.Since(initTime)
	fmt.Println("Time elasped: ", elapsedTime)
	output_message_files = append(output_message_files, "Time elapsed: " + elapsedTime.String())
	output_message_err = append(output_message_err, "Time elapsed: " + elapsedTime.String())

	if len(output_message_err) == 0{
		fmt.Println("No Errors")
	}else{
		file, err := os.OpenFile("./output/errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		err = os.Truncate("./output/errors.txt", 0)

		if err != nil {
			log.Fatalf("failed truncating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		datawriter.WriteString(time.Now().String() + "\n")

		for _, data := range output_message_err {
			_, _ = datawriter.WriteString(data + "\n")
		}
	 
		datawriter.Flush()
		file.Close()

	}


	if len(output_message_files) == 0{
		fmt.Println("No Results")
	}else{
		file, err := os.OpenFile("./output/results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		err = os.Truncate("./output/results.txt", 0)

		if err != nil {
			log.Fatalf("failed truncating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		datawriter.WriteString(time.Now().String() + "\n")

		for _, data := range output_message_files {
			_, _ = datawriter.WriteString(data + "\n")
		}
	 
		datawriter.Flush()
		file.Close()

	}
}

func startFileSearch(keyword string, pwd string, count int, depth int){


	if configuration.SearchType == 0{

		//deep activation
		wg_deep.Add(1)
		
		deepSearchFile(keyword, pwd, count, depth)

		wg_deep.Wait()

	}else if configuration.SearchType == 1{

		//shallow activation
		stack = append(stack, pwd)

		for len(stack) != 0{
			dir := stack[0]
			stack = stack[1:]
			shallowSearchFile(keyword, dir, count) //allow for manual input of shallow vs deep
		}
		shallowDepth++

		for len(newStack) != 0{
			stack = make([]string, len(newStack))
			copy(stack, newStack)
			newStack = []string{}
			for len(stack) != 0{
				dir := stack[0]
				stack = stack[1:]
				shallowSearchFile(keyword, dir, count) //allow for manual input of shallow vs deep
			}
			shallowDepth++
		}
	}
	cGUI <- "END"

	wg_main.Done()
}

func callCLIGUI(){
	payload := ""
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

	wg_main.Done()
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
	fmt.Print("Enter keyword to search for: ")
	fmt.Scanln(&keyword)

	return keyword
}

func logErr(err error){
	message := strings.Split(err.Error(), ": ")
	if message[len(message)-1] == "file name too long"{
		output_message_err = append(output_message_err, err.Error())
	}else if message[len(message)-1] == "no such file or directory"{
		output_message_err = append(output_message_err, err.Error())
	}else if message[len(message)-1] == "permission denied"{
		output_message_err = append(output_message_err, err.Error())
	}else if message[len(message)-1] == "operation not permitted"{
		output_message_err = append(output_message_err, err.Error())
	}else if message[len(message)-1] == "bad file descriptor"{
		output_message_err = append(output_message_err, err.Error())
	}else{
		log.Fatal(err)
	}
}

func deepSearchFile(keyword string, pwd string, count int, depth int){
	
	files, err := os.ReadDir(pwd)

	if err != nil {
		logErr(err)
		wg_deep.Done()
		return
	}

	for _, f := range files {

		if f.Name()[0] == 46{ //dot
			continue
		}

		count++
		extraSlash := ""

		if pwd != "/"{
			extraSlash = "/"
		}

		fstats, err := os.Stat(pwd + extraSlash + f.Name())

		if err != nil {
			logErr(err)
			continue

		}else{

			curr_pwd := pwd + (extraSlash + fstats.Name())

			if configuration.ExactMatch{

				if keyword == fstats.Name(){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					output_message_files = append(output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}

			}else if !configuration.ExactMatch{

				if strings.Contains(fstats.Name(), keyword){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					output_message_files = append(output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}
				
			}
			
			if fstats.IsDir(){

				if depth >= configuration.MAX_DEPTH-1{
					continue
				}
				wg_deep.Add(1)
				go deepSearchFile(keyword, curr_pwd, count, depth+1)
			}
		}
	}
	wg_deep.Done()
}

func shallowSearchFile(keyword string, pwd string, count int){
	
	files, err := os.ReadDir(pwd)
	
	if err != nil {
		logErr(err)
		return
	}

	for _, f := range files {

		if f.Name()[0] == 46{ //dot
			continue
		}

		count++
		extraSlash := ""

		if pwd != "/"{
			extraSlash = "/"
		}

		fstats, err := os.Stat(pwd + extraSlash + f.Name())

		if err != nil {
			logErr(err)
			continue

		}else{

			curr_pwd := pwd + (extraSlash + fstats.Name())

			if configuration.ExactMatch{

				if keyword == fstats.Name(){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					output_message_files = append(output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}

			}else if !configuration.ExactMatch{

				if strings.Contains(fstats.Name(), keyword){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					output_message_files = append(output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}

			}
			
			if fstats.IsDir(){

				if shallowDepth >= configuration.MAX_DEPTH-1{
					continue
				}
				newStack = append(newStack, curr_pwd)
			}
		}

	}
}




