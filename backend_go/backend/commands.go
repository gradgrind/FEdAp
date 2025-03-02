package backend

import (
	"gradgrind/backend/base"
	"gradgrind/backend/fet"
	"gradgrind/backend/ttbase"
	"gradgrind/backend/w365tt"
	"strconv"
	"time"
)

var DB *base.DbTopLevel   // The current database
var TtData *ttbase.TtInfo // The current timetable data

var commandMap map[string]func(map[string]any, map[string]any) string

func init() {
	commandMap = map[string]func(map[string]any, map[string]any) string{}

	commandMap["SLEEP"] = testSleep

	commandMap["SET_LANGUAGE"] = setLanguage
	commandMap["LOAD_W365_JSON"] = loadW365Json
	commandMap["SAVE_FET"] = makeFetFiles
}

// commandHandler is the dispatcher loop, reading commands from the input
// channel (xchan) and calling corresponding handlers. Responses are written
// to the output channel (ochan).
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

func ReportCancelled(cmd map[string]any, outmap map[string]any) string {
	base.Report(`<Notice>Operation cancelled>`)
	outmap["DATA"] = cmd
	return "CANCELLED"
}

// This command is just for testing
func testSleep(cmd map[string]any, outmap map[string]any) string {
	tsecs := int(cmd["TIME"].(float64))
	for i := range tsecs {
		if cancel {
			return ReportCancelled(cmd, outmap)
		}
		time.Sleep(1 * time.Second)
		base.Report(`<Info>Tick>`)
		base.Report(`<PROGRESS>%s>`, strconv.Itoa(i+1))
	}
	outmap["TIME"] = tsecs
	return "SLEPT"
}

// TODO: setLanguage reads a translations file ...
func setLanguage(cmd map[string]any, outmap map[string]any) string {
	if base.ReadMessages("/home/user/tmp/messages") {
		return "OK"
	}
	return "FAILED"
}

// loadW365Json reads a Waldorf 365 timetable-data file (_w365.json) and
// sets up the data as the current database.
func loadW365Json(cmd map[string]any, outmap map[string]any) string {
	//TODO-- Start testing with a single fixed file?
	cmd["FILEPATH"] = "/home/user/tmp/test1_w365.json"

	if LoadFile(cmd, w365tt.LoadJSON) {
		return "OK"
	} else {
		return "FAILED"
	}
}

// makeFetFiles generates a FET file (.fet) from the current database.
func makeFetFiles(cmd map[string]any, outmap map[string]any) string {
	var fetfile string
	var mapfile string

	/*TODO++
		fetfile = cmd["OUTFILE"].(string) // without ending!
	    mf, ok := cmd["MAPFILE"]
	    if ok {
	        mapfile = mf.(string)
	    } else {
	        mapfile = fetfile
	    }
	*/

	//TODO--
	fetfile = "/home/user/tmp/test1"
	mapfile = fetfile

	if TtData == nil {
		TtData = ttbase.MakeTtInfo(DB)
	}
	if TtData.Placements == nil {
		TtData.PrepareCoreData()
	}
	xmlitem, lessonIdMap := fet.MakeFetFile(TtData)

	if fetfile != "" {
		fetfile += ".fet"
		if !SaveFile(fetfile, []byte(xmlitem)) {
			return "FAILED"
		}
		base.Report(`<Info>FET file written to: %s>`, fetfile)
	}
	if mapfile != "" {
		mapfile += ".map"
		if !SaveFile(mapfile, []byte(lessonIdMap)) {
			return "FAILED"
		}
		base.Report(`<Info>Id-map written to: %s>`, mapfile)
	}
	return "OK"
}
