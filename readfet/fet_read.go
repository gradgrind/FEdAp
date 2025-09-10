package readfet

import (
	"encoding/xml"
	"io"
	"os"
)

func ReadFet(fetpath string) (*Fet, error) {
	// Open the  XML file
	xmlFile, err := os.Open(fetpath)
	if err != nil {
		return nil, err
	}
	// Remember to close the file at the end of the function
	defer xmlFile.Close()
	// Read the opened XML file as a byte array.
	//fmt.Printf("Reading: %s\n", fetpath)
	byteValue, _ := io.ReadAll(xmlFile)

	// Parse XML to FET structure
	fetdata := &Fet{}
	xml.Unmarshal(byteValue, fetdata)
	return fetdata, nil
}
