package readxml

import (
	"fmt"
	"gradgrind/backend/base"
	"gradgrind/backend/fet"
	"gradgrind/backend/ttbase"
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
	base.OpenLog(nil, "")
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
		base.Report(`<Error>Opening FET file: %s>`, err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(xmlitem)
	if err != nil {
		base.Report(`<Error>Writing FET file: %s>`, err)
		return
	}
	base.Report(`<Info>Wrote FET file to: %s>`, fetfile)

	// Write Id-map file.
	mapfile := fetpath + ".map"
	fm, err := os.Create(mapfile)
	if err != nil {
		base.Report(`<Error>Opening map file: %s>`, err)
		return
	}
	defer fm.Close()
	_, err = fm.WriteString(lessonIdMap)
	if err != nil {
		base.Report(`<Error>Writing map file: %s>`, err)
		return
	}
	base.Report(`<Info>Wrote Id-map to: %s>`, mapfile)
}
