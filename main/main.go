package main

import (
	"fmt"
	"os"
	//"os/exec"
	//"strconv"
	"encoding/json"
	"strings"
	"sync"
	"time"
	"runtime"
	"github.com/atdykema/file_explorer_go/imageSearch"
	"github.com/atdykema/file_explorer_go/logging"
	"github.com/atdykema/file_explorer_go/testing"
	"path/filepath"
)

type Configuration struct {
	CommandLineConfig bool
	MAX_DEPTH int
	SearchType int
	ExactMatch bool
	PrintConfigOnSearch bool
	PrintResults bool
	ImageEnabled bool
}

//!!!!!!!!!!!!!!!!!!!!

//TODO: allow search for file type, last modification

var cGUI chan string
var wg_main *sync.WaitGroup
var wg_deep *sync.WaitGroup
var stack []string
var newStack []string
var shallowDepth int
var configuration Configuration
//var cc int = 0
//var mutex_search sync.Mutex

func initConfig() Configuration{
	absPath, _ := filepath.Abs("../config/config.json")
	file, _ := os.Open(absPath)
	
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	
	err := decoder.Decode(&configuration)
	
	if err != nil {	
		logging.LogErr(err)
	}
	file.Close()
	return configuration
}

func init(){
	//init config
	
	if configuration.ImageEnabled{
		imageSearch.InitImg()
	}

	configuration = initConfig()
	
	cGUI = make(chan string)
	wg_main = &sync.WaitGroup{}
	wg_deep = &sync.WaitGroup{}
	stack = []string{}
	newStack = []string{}
	shallowDepth = 0

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main(){

	testing.Main_test()

	if configuration.PrintConfigOnSearch{
		fmt.Println("Config:")
		fmt.Println("\tCommandLineConfig: ", configuration.CommandLineConfig)
		fmt.Println("\tMAX_DEPTH: ", configuration.MAX_DEPTH)
		fmt.Println("\tSearchType: ", configuration.SearchType)
		fmt.Println("\tExactMatch: ", configuration.ExactMatch)
		fmt.Println("\tPrintConfigOnSearch: ", configuration.PrintConfigOnSearch)
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
	logging.Output_message_files = append(logging.Output_message_files, "Time elapsed: " + elapsedTime.String())
	logging.Output_message_err = append(logging.Output_message_err, "Time elapsed: " + elapsedTime.String())

	logging.LogOutput()
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

		if configuration.PrintResults{
			fmt.Println(payload)
		}
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

func deepSearchFile(keyword string, pwd string, count int, depth int){

	files, err := os.ReadDir(pwd)

	if err != nil {
		logging.LogErr(err)
		/*
		cc--
		fmt.Println(cc)
		*/
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
			logging.LogErr(err)
			continue

		}else{

			curr_pwd := pwd + (extraSlash + fstats.Name())

			if configuration.ExactMatch{

				if keyword == fstats.Name(){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					logging.Output_message_files = append(logging.Output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}
				

			}else if !configuration.ExactMatch{

				if strings.Contains(fstats.Name(), keyword){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					logging.Output_message_files = append(logging.Output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}
				
			}
			if fstats.IsDir(){

				if depth >= configuration.MAX_DEPTH-1{
					continue
				}
				wg_deep.Add(1)
				/*
				cc++
				fmt.Println(cc)
				*/
				go deepSearchFile(keyword, curr_pwd, count, depth+1)
			}
		}
	}
	/*
	cc--
	fmt.Println(cc)
	*/
	wg_deep.Done()
}

func shallowSearchFile(keyword string, pwd string, count int){
	
	files, err := os.ReadDir(pwd)
	
	if err != nil {
		logging.LogErr(err)
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
			logging.LogErr(err)
			continue

		}else{

			curr_pwd := pwd + (extraSlash + fstats.Name())

			if configuration.ExactMatch{

				if keyword == fstats.Name(){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					logging.Output_message_files = append(logging.Output_message_files, curr_pwd)
					cGUI <- curr_pwd
				}

			}else if !configuration.ExactMatch{

				if strings.Contains(fstats.Name(), keyword){
					//output_message_files[0] = ("Files searched: " + strconv.Itoa(count))
					logging.Output_message_files = append(logging.Output_message_files, curr_pwd)
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




