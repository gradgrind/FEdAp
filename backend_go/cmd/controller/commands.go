package main

import (
	"gradgrind/backend/base"
	"gradgrind/backend/w365tt"
	"strconv"
	"time"
)

var commandMap map[string]func(map[string]any, map[string]any) string

func init() {
	commandMap = map[string]func(map[string]any, map[string]any) string{}

	commandMap["SLEEP"] = testSleep

	commandMap["LOAD_W365_JSON"] = loadW365Json
}

// TODO: Not working yet ...
func commandHandler(ochan chan map[string]any, xchan chan map[string]any) {
	var done string
	for {
		odata := map[string]any{}
		xdata := <-xchan
		running = true
		cmd, ok := xdata["DO"].(string)
		if ok {
			if f, ok := commandMap[cmd]; ok {
				done = f(xdata, odata)
				goto done_send
			}
		}
		done = "UNKNOWN_COMMAND"
		odata["DATA"] = xdata

	done_send:

		odata["DONE"] = done
		ochan <- odata

		running = false
	}
}

// TODO: Move to base?
func ReportCancelled() string {
	base.Report(`<Notice>Operation cancelled>`)
	return "CANCELLED"
}

// This command is just for testing
func testSleep(cmd map[string]any, outmap map[string]any) string {
	tsecs := int(cmd["TIME"].(float64))
	for i := range tsecs {
		if cancel {
			outmap["DATA"] = cmd
			return ReportCancelled()
		}
		time.Sleep(1 * time.Second)
		base.Report(`<Info>Tick>`)
		base.Report(`<PROGRESS>%s>`, strconv.Itoa(i+1))
	}
	outmap["TIME"] = tsecs
	return "SLEPT"
}

func loadW365Json(cmd map[string]any, outmap map[string]any) string {
	//TODO: Maybe start testing with a single fixed file?
	fpath := "/home/user/tmp/test1_w365.json"
	w365tt.LoadFile(fpath)
	return "LOADED"
}
