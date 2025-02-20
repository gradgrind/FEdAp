package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	unsaved := true

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
				if unsaved {
					if idata["FORCE"] != true {
						odata = map[string]any{
							"DONE": "QUIT_UNSAVED",
						}
						ochan <- odata
						continue
					}
				}
				os.Exit(0)
			}
			//TODO...

			odata = map[string]any{
				"DONE": idata,
			}
			ochan <- odata
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

func sender(ochan chan map[string]any, writer *bufio.Writer) {
	data := <-ochan
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal json: %s\n", err))
	}

	fmt.Fprintln(writer, string(jsonData))
	writer.Flush() // Don't forget to flush!
}
