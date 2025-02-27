package main

import (
	"gradgrind/backend/w365tt"
	"strconv"
	"time"
)

var commandMap map[string]func(map[string]any, map[string]any) string

func init() {
	commandMap = map[string]func(map[string]any, map[string]any) string{}

	commandMap["LOAD_W365_JSON"] = loadW365Json
}

func commandHandler(ochan chan map[string]any, xchan chan map[string]any) {
	for {
		odata := map[string]any{}
		var done string
		xdata := <-xchan
		running = true

		//TODO: If there's only SLEEP and default, a switch is not really
		// appropriate!
		// Also, if I return the done-map, I don't need to pass it as argument.
		// On the other hand, it might be better to create it once and
		// clear it on each repeat, rather than building a new one each time?

		switch cmd := xdata["DO"].(string); cmd {
		case "SLEEP": // for testing
			tsecs := int(xdata["TIME"].(float64))
			for i := range tsecs {
				if cancel {
					ochan <- map[string]any{
						"DONE":   "",
						"REPORT": "Notice",
						"TEXT":   "OPERATION_CANCELLED",
					}
					done = "CANCELLED"
					odata["DATA"] = xdata
					goto send_done
				}
				time.Sleep(1 * time.Second)
				ochan <- map[string]any{
					"DONE":   "",
					"REPORT": "Info",
					"TEXT":   "TICK",
				}
				ochan <- map[string]any{
					"DONE":   "",
					"REPORT": "PROGRESS",
					"TEXT":   strconv.Itoa(i + 1),
				}
			}
			done = "SLEPT"
			odata["TIME"] = tsecs

		default:
			f, ok := commandMap[cmd]
			if ok {
				done = f(xdata, odata)
			} else {
				done = "UNKNOWN_COMMAND"
				odata["DATA"] = xdata
			}
		}
	send_done:
		odata["DONE"] = done
		ochan <- odata
		running = false
	}
}

func loadW365Json(cmd map[string]any, outmap map[string]any) string {
	//TODO: Maybe start testing with a single fixed file?
	fpath := "/home/user/tmp/test1_w365.json"
	w365tt.LoadFile(fpath)
	return "LOADED"
}
