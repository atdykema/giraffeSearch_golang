package logging

import(
	"strings"
	"fmt"
	"path/filepath"
	"os"
	"log"
	"bufio"
	"time"
)

var Output_message_files []string
var Output_message_err []string

func init(){
	Output_message_files = make([]string, 0)
	Output_message_err = make([]string, 0)
}

func LogErr(err error){
	message := strings.Split(err.Error(), ": ")
	if message[len(message)-1] == "file name too long"{
		Output_message_err = append(Output_message_err, err.Error())
	}else if message[len(message)-1] == "no such file or directory"{
		Output_message_err = append(Output_message_err, err.Error())
	}else if message[len(message)-1] == "permission denied"{
		Output_message_err = append(Output_message_err, err.Error())
	}else if message[len(message)-1] == "operation not permitted"{
		Output_message_err = append(Output_message_err, err.Error())
	}else if message[len(message)-1] == "bad file descriptor"{
		Output_message_err = append(Output_message_err, err.Error())
	}else{
		Output_message_err = append(Output_message_err, err.Error())
	}
}

func LogOutput(){

	if len(Output_message_err) == 0{
		fmt.Println("No Errors")
	}else{
		absPath, _ := filepath.Abs("../output/errors.txt")
		file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		err = os.Truncate(absPath, 0)

		if err != nil {
			log.Fatalf("failed truncating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		datawriter.WriteString(time.Now().String() + "\n")

		for _, data := range Output_message_err {
			_, _ = datawriter.WriteString(data + "\n")
		}
	 
		datawriter.Flush()
		file.Close()

	}


	if len(Output_message_files) == 0{
		fmt.Println("No Results")
	}else{
		absPath, _ := filepath.Abs("../output/results.txt")
		file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		err = os.Truncate(absPath, 0)

		if err != nil {
			log.Fatalf("failed truncating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		datawriter.WriteString(time.Now().String() + "\n")

		for _, data := range Output_message_files {
			_, _ = datawriter.WriteString(data + "\n")
		}
	 
		datawriter.Flush()
		file.Close()

	}
}