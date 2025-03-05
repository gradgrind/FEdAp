package backend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gradgrind/backend/base"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	// WorkingDir is the directory into which output files will be written
	WorkingDir string
	// StemFile is the name of the current data set, the data-file name
	// without type suffix
	StemFile string
	// LogFile is the path relative to WorkingDir of the log file for the
	// current run
	LogFile string
	// SenderChannel is used to transmit messages
	SenderChannel chan map[string]any
)

func init() {
	//TODO

	WorkingDir = "/home/user/tmp/FEdAp"

	expath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exname := filepath.Base(expath)
	exstem := strings.TrimSuffix(exname, filepath.Ext(exname))
	LogFile = filepath.Base(exstem) + ".log"

}

// It should be possible to cancel long-running operations.
// However, a goroutine can't be stopped from the outside, so any
// long-running operation should periodically check whether the [cancel]
// flag is set (true).
var cancel bool

// When an operation is running, no new operations can be started, the
// [running] flag is set (true). It should be regarded as an error if a
// new comman request arrives while this flag is set. To ensure that
// the front-end knows when an operation has completed, all completions
// should send a "DONE" message.
var running bool

func BackEnd() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	//TODO: Maybe unsaved should be a global variable?
	unsaved := rand.IntN(3) == 0

	SenderChannel = make(chan map[string]any)
	go sender(writer)
	xchan := make(chan map[string]any)
	go commandHandler(xchan)

	logpath := filepath.Join(WorkingDir, LogFile)
	base.OpenLog(SenderChannel, logpath)

	for {
		message, _ := reader.ReadString('\n')

		var idata map[string]any
		var odata map[string]any
		err := json.Unmarshal([]byte(message), &idata)
		if err != nil {
			e := fmt.Sprintf("Could not unmarshal json: %s\n:: %s\n",
				err, message)
			odata = map[string]any{
				"DONE": "ERROR",
				"TEXT": e,
			}
			SenderChannel <- odata
			continue
		} else if n, ok := idata["DO"]; ok {
			if n == "QUIT" {
				// If there are unsaved changes, respond with "DONE": "QUIT_UNSAVED"
				if unsaved && idata["FORCE"] != true {
					odata = map[string]any{
						"DONE":   "",
						"REPORT": "QUIT_UNSAVED?",
					}
					SenderChannel <- odata
					continue
				}
				os.Exit(0)
			}
			if n == "CANCEL" {
				if running {
					cancel = true
				}
				continue
			}
			// Pass command to handler.
			// Ensure somehow that only one gets handled at a time.
			if running {
				odata = map[string]any{
					"DONE":   "",
					"REPORT": "BACKEND_BUSY",
					"DATA":   idata,
				}
				SenderChannel <- odata
				continue
			}
			cancel = false
			xchan <- idata
			continue
		}
		//TODO: No "DO" ...

		odata = map[string]any{
			"DONE": "ERROR",
			"DATA": idata,
		}
		SenderChannel <- odata
	}
}

// TODO: Maybe this needs the possibility to block sending?
func sender(writer *bufio.Writer) {
	for {
		data := <-SenderChannel
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(fmt.Sprintf("Could not marshal json: %s\n", err))
		}

		fmt.Fprintln(writer, string(jsonData))
		writer.Flush() // Don't forget to flush!
	}
}
