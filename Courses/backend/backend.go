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
		} else if n, ok := idata["QUIT"]; ok {
			nn := int(n.(float64))
			if nn == 0 {
				odata = idata
			} else {
				os.Exit(nn)
			}
		} else {
			odata = map[string]any{
				"GOT": idata,
			}
		}

		jsonData, err := json.Marshal(odata)
		if err != nil {
			panic(fmt.Sprintf("Could not marshal json: %s\n", err))
		}

		fmt.Fprintln(writer, string(jsonData))
		writer.Flush() // Don't forget to flush!
	}
}
