package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	//TODO:
	unsaved := rand.IntN(3) == 0

	ochan := make(chan map[string]any)
	go sender(ochan, writer)

	for {
		message, _ := reader.ReadString('\n')

		var idata map[string]any
		var odata map[string]any
		err := json.Unmarshal([]byte(message), &idata)
		if err != nil {
			e := fmt.Sprintf("Could not unmarshal json: %s\n:: %s\n",
				err, message)
			odata = map[string]any{
				"ERROR": e,
			}
		} else if n, ok := idata["DO"]; ok {
			if n == "QUIT" {
				// If there are unsaved changes, respond with "DONE": "QUIT_UNSAVED"
				if unsaved && idata["FORCE"] != true {
					odata = map[string]any{
						"DONE":   "",
						"REPORT": "QUIT_UNSAVED?",
					}
					ochan <- odata
					continue
				}
				os.Exit(0)
			}
			//TODO...

			// This one sets up a new goroutine for each command.
			// That might not be a bad idea, but I should ensure somehow that
			// only one gets handled at a time ... in which case I could
			// use just one goroutine with a read loop.
			go handle_command(ochan, idata)
			continue
		}
		//TODO: No "DO" ...

		odata = map[string]any{
			"DONE": "ERROR",
			"DATA": idata,
		}
		ochan <- odata
	}
}

// TODO: Maybe this needs the possibility to block sending?
func sender(ochan chan map[string]any, writer *bufio.Writer) {
	for {
		data := <-ochan
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(fmt.Sprintf("Could not marshal json: %s\n", err))
		}

		fmt.Fprintln(writer, string(jsonData))
		writer.Flush() // Don't forget to flush!
	}
}
