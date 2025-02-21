package main

import "time"

func commandHandler(ochan chan map[string]any, xchan chan map[string]any) {
	for {
		odata := map[string]any{}
		var done string
		xdata := <-xchan
		running = true

		switch cmd := xdata["DO"]; cmd {
		case "SLEEP": // for testing
			tsecs := int(xdata["TIME"].(float64))
			for range tsecs {
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
			}
			done = "SLEPT"
			odata["TIME"] = tsecs
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
