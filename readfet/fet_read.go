package readfet

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func ReadFet(fetpath string) *Fet {

	fmt.Println("\n ++++++++++++++++++++++")

	// Open the  XML file
	xmlFile, err := os.Open(fetpath)
	if err != nil {
		panic(err)
	}
	// Remember to close the file at the end of the function
	defer xmlFile.Close()
	// read the opened XML file as a byte array.
	fmt.Printf("Reading: %s\n", fetpath)
	byteValue, _ := io.ReadAll(xmlFile)

	// Parse XML to FET structure
	fetdata := &Fet{}
	xml.Unmarshal(byteValue, fetdata)
	return fetdata
}
