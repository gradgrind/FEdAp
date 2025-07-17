package main

import (
	"fedap/base"
	"fedap/fet"
	"fmt"
	"io"
	"os"
)

var inputfiles = []string{
	"../../testdata/readxml/Demo1.fet",
	"../../testdata/readxml/x01.fet",
}

func main() {

	base.OpenLog("")
	for _, fetpath := range inputfiles {
		fmt.Println("\n ++++++++++++++++++++++")

		// Open the  XML file
		xmlFile, err := os.Open(fetpath)
		if err != nil {
			base.Error.Fatal(err)
		}
		// Remember to close the file at the end of the function
		defer xmlFile.Close()
		// read the opened XML file as a byte array.
		base.Message.Printf("Reading: %s\n", fetpath)
		byteValue, _ := io.ReadAll(xmlFile)

		// Parse XML to FET structure
		fetdata := fet.ReadFet(byteValue)

		tt_data := &fet.TtData{}
		tt_data.PrepareResources(fetdata)
		tt_data.SetupActivities(fetdata)
		tt_data.ResourceBlocking(fetdata)
		tt_data.SetupFixedTimes(fetdata)
		tt_data.SetupDaysBetween(fetdata)
		tt_data.ConnectBags()

		//TODO--
		//tt_data.PrintBags()

		//TODO: place fixed activities, generate sets of placement groups

		/*
			db := cdata.Db()
			db.PrepareDb()
			ttbase.MakeTtInfo(db)
			ttinfo := ttbase.MakeTtInfo(db)
			ttinfo.PrepareCoreData()

			j := ttbase.TtInfoToJson(ttinfo)
			jfp := filepath.Base(fxml)
			jfp = strings.TrimSuffix(jfp, filepath.Ext(jfp)) + "_tt.json"
			err := os.WriteFile(jfp, j, 0644)
			if err != nil {
				panic(err)
			}
		*/

		//ttengine.PlaceLessons(ttinfo)
	}
}
