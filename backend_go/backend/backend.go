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

	ochan := make(chan map[string]any)
	go sender(ochan, writer)
	xchan := make(chan map[string]any)
	go commandHandler(ochan, xchan)

	logpath := filepath.Join(WorkingDir, LogFile)
	base.OpenLog(ochan, logpath)

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
				ochan <- odata
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
