package readxml

import (
	"fedap/base"
	"fedap/fet"
	"fedap/ttbase"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var inputfiles = []string{
	"../testdata/readxml/Demo1.xml",
	"../testdata/readxml/x01.xml",
}

func Test2JSON(t *testing.T) {
	base.OpenLog("")
	for _, fxml := range inputfiles {
		fmt.Println("\n ++++++++++++++++++++++")
		cdata := ConvertToDb(fxml)
		fmt.Println("*** Available Schedules:")
		slist := cdata.ScheduleNames()
		for _, sname := range slist {
			fmt.Printf("  -- %s\n", sname)
		}
		sname := "Vorlage"
		if !slices.Contains(slist, sname) {
			if len(slist) != 0 {
				sname = slist[0]
			} else {
				fmt.Println(" ... stopping ...")
				continue
			}
		}
		fmt.Printf("*** Using Schedule '%s'\n", sname)
		if !cdata.ReadSchedule(sname) {
			fmt.Println(" ... failed ...")
			continue
		}
		stempath := strings.TrimSuffix(fxml, filepath.Ext(fxml))
		fjson := stempath + "_db.json"
		if cdata.db.SaveDb(fjson) {
			fmt.Printf("\n ***** Written to: %s\n", fjson)
		} else {
			fmt.Println("\n ***** Write to JSON failed")
			continue
		}

		stempath = strings.TrimSuffix(stempath, "_w365")
		toFET(cdata.db, stempath)
	}
}

func toFET(db *base.DbTopLevel, fetpath string) {
	db.PrepareDb()
	ttinfo := ttbase.MakeTtInfo(db)
	ttinfo.PrepareCoreData()

	// ********** Build the fet file **********
	xmlitem, lessonIdMap := fet.MakeFetFile(ttinfo)

	// Write FET file
	fetfile := fetpath + ".fet"
	f, err := os.Create(fetfile)
	if err != nil {
		base.Bug.Fatalf("Couldn't open output file: %s\n", fetfile)
	}
	defer f.Close()
	_, err = f.WriteString(xmlitem)
	if err != nil {
		base.Bug.Fatalf("Couldn't write fet output to: %s\n", fetfile)
	}
	base.Message.Printf("FET file written to: %s\n", fetfile)

	// Write Id-map file.
	mapfile := fetpath + ".map"
	fm, err := os.Create(mapfile)
	if err != nil {
		base.Bug.Fatalf("Couldn't open output file: %s\n", mapfile)
	}
	defer fm.Close()
	_, err = fm.WriteString(lessonIdMap)
	if err != nil {
		base.Bug.Fatalf("Couldn't write fet output to: %s\n", mapfile)
	}
	base.Message.Printf("Id-map written to: %s\n", mapfile)
	base.Message.Println("OK")
}
