package main

import (
	"gradgrind/backend/w365tt"
	"strconv"
	"time"
)

func commandHandler(ochan chan map[string]any, xchan chan map[string]any) {
	for {
		odata := map[string]any{}
		var done string
		xdata := <-xchan
		running = true

		switch cmd := xdata["DO"]; cmd {
		case "SLEEP": // for testing
			tsecs := int(xdata["TIME"].(float64))
			for i := range tsecs {
				if cancel {
					ochan <- map[string]any{
						"DONE":   "",
						"REPORT": "REPORT",
						"TEXT":   "OPERATION_CANCELLED",
					}
					done = "CANCELLED"
					odata["DATA"] = xdata
					goto send_done
				}
				time.Sleep(1 * time.Second)
				ochan <- map[string]any{
					"DONE":   "",
					"REPORT": "REPORT",
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

		case "LOAD_W365_JSON":
			//TODO: Maybe start testing with a single fixed file?
			fpath := "/home/user/tmp/test1_w365.json"
			w365tt.LoadFile(ochan, fpath)
			done = "LOADED"

		default:
			done = "UNKNOWN_COMMAND"
			odata["DATA"] = xdata
		}
	send_done:
		odata["DONE"] = done
		ochan <- odata
		running = false
	}
}
