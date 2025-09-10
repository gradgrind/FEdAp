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
	if err != nil {
		panic(err)
	}

	fetdata := readfet.ReadFet(abspath)

	// Generate output
	stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	fet2xlsx.GetTeachers(fetdata)
	fet2xlsx.GetStudentGroups(fetdata)
	activities := fet2xlsx.GetActivityData(fetdata)
	opath, err := fet2xlsx.TeachersActivities(fetdata, activities, stempath)
	if err == nil {
		zenity.Info("Generated: "+opath,
			zenity.Title("Information"),
			zenity.InfoIcon)
	} else {
		zenity.Error(err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon)
	}
	opath, err = fet2xlsx.StudentsActivities(fetdata, activities, stempath)
	if err == nil {
		zenity.Info("Generated: "+opath,
			zenity.Title("Information"),
			zenity.InfoIcon)
	} else {
		zenity.Error(err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon)
	}
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
