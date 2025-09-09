package main

import (
	"fmt"
	"strings"

	"github.com/ncruces/zenity"

	"fedap/fet2xlsx"
	"fedap/readfet"
	"flag"
	"path/filepath"
)

//TODO: Replace panic calls by something to show the message

const defaultPath = ``

func main() {
	abspath, err := zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{
				Name:     "FET result files",
				Patterns: []string{"*_data_and_timetable.fet"},
				CaseFold: false,
			},
		})
	if err == nil {
		fmt.Printf("Reading %s\n", abspath)
	} else {
		panic(err)
	}

	fetdata := readfet.ReadFet(abspath)

	//TODO: generate output
	//fmt.Printf(" --->\n%v\n", fetdata)

	stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	activities := fet2xlsx.GetActivityData(fetdata)
	fet2xlsx.TeachersActivities(fetdata, activities, stempath)
}

func maincli() {

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		if len(args) == 0 {
			panic("ERROR* No input file")
		}
		panic(fmt.Sprintf("*ERROR* Too many command-line arguments:\n  %+v\n", args))
	}
	abspath, err := filepath.Abs(args[0])
	if err != nil {
		panic(fmt.Sprintf("*ERROR* Couldn't resolve file path: %s\n", args[0]))
	}

	fetdata := readfet.ReadFet(abspath)

	//TODO: generate output
	//fmt.Printf(" --->\n%v\n", fetdata)

	//stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	fet2xlsx.GetActivityData(fetdata)
	//fet2xlsx.TeachersActivities(fetdata)
}
