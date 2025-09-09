package main

import (
	"fedap/fet2xlsx"
	"fedap/readfet"
	"flag"
	"log"
	"path/filepath"
)

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		if len(args) == 0 {
			log.Fatalln("ERROR* No input file")
		}
		log.Fatalf("*ERROR* Too many command-line arguments:\n  %+v\n", args)
	}
	abspath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatalf("*ERROR* Couldn't resolve file path: %s\n", args[0])
	}

	fetdata := readfet.ReadFet(abspath)

	//TODO: generate output
	//fmt.Printf(" --->\n%v\n", fetdata)

	//stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	fet2xlsx.GetActivityData(fetdata)
	//fet2xlsx.TeachersActivities(fetdata)
}
