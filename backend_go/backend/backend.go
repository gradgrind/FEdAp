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

type Ref = base.Ref

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
	// It should be possible to cancel long-running operations.
	// However, a goroutine can't be stopped from the outside, so any
	// long-running operation should periodically check whether the [cancel]
	// flag is set (true).
	cancel bool
	// When an operation is running, no new operations can be started, the
	// [running] flag is set (true). It should be regarded as an error if a
	// new comman request arrives while this flag is set. To ensure that
	// the front-end knows when an operation has completed, all completions
	// should send a "DONE" message.
	running bool
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
		}
		if n, ok := idata["DO"]; ok {
			if n == "QUIT" {
				// If there are unsaved changes, request a "QUIT_UNSAVED?"
				// dialog.
				if unsaved && idata["FORCE"] != true {
					odata = map[string]any{
						"DONE":   false,
						"DIALOG": "QUIT_UNSAVED?",
						"TEXT": base.Tr(
							`<>There is unsaved data. Quit and lose it?>`),
					}
					SenderChannel <- odata
					continue
				}
				odata = map[string]any{
					"DONE": true,
				}
				SenderChannel <- odata
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
				base.Report(`<Bug>Unexpected front-end message (back-end
					> busy):\n -- %s>`, message)
				//TODO: Should this be sent?
				odata = map[string]any{
					"DONE": false,
				}
				SenderChannel <- odata
				continue
			}
			// Pass the message to the command handler.
			cancel = false
			xchan <- idata
			continue
		}
		// No "DO" ...
		base.Report(`<Bug>Invalid front-end message:\n -- %s>`, message)
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

func gui(cmd string, object string, data any) {
	payload := map[string]any{
		"GUI":    cmd,
		"OBJECT": object,
		"DATA":   data,
	}
	SenderChannel <- payload
}

func getElementNames(ref Ref) []string {
	e, ok := DB.Elements[ref]
	if !ok {
		base.Report("<Bug>Unknown Element: %s>", ref)
		return []string{"", ""}
	}
	en := e.GetElementStrings()
	return []string{en.Short, en.Long}
}
