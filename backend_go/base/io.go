package base

import (
	"encoding/json"
	"io"
	"os"
)

func init() {
	Tr(map[string]string{
		"BUG_SAVE_JSON":          "[Bug] Saving JSON: %v",
		"ERROR_SAVE_FILE":        "[Error] Saving file: %v",
		"ERROR_OPEN_FILE":        "[Error] Opening file: %v",
		"INFO_READ_FILE":         "[Info] Reading file: %s",
		"ERROR_BAD_JSON":         "[Error] Invalid JSON: %s",
		"ERROR_ELEMENT_NO_ID":    "[Error] Element has no Id:\n  -- %+v",
		"ERROR_ELEMENT_ID_REUSE": "[Error] Element Id defined more than once:\n  %s",
	})
}

func (db *DbTopLevel) SaveDb(fpath string) bool {
	// Save as JSON
	j, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		Report("BUG_SAVE_JSON", err)
		return false
	}
	if err := os.WriteFile(fpath, j, 0666); err != nil {
		Report("ERROR_SAVE_FILE", err)
		return false
	}
	return true
}

func LoadDb(fpath string) *DbTopLevel {
	// Open the  JSON file
	jsonFile, err := os.Open(fpath)
	if err != nil {
		Report("ERROR_OPEN_FILE", err)
		return nil
	}
	// Remember to close the file at the end of the function
	defer jsonFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	Report("INFO_READ_FILE", fpath)
	v := NewDb()
	err = json.Unmarshal(byteValue, v)
	if err != nil {
		Report("ERROR_BAD_JSON", err)
		return nil
	}
	v.initElements()
	return v
}

func (db *DbTopLevel) testElement(ref Ref, element Elem) {
	if ref == "" {
		Report("ERROR_ELEMENT_NO_ID", element)
		db.SetInvalid()
	}
	_, nok := db.Elements[ref]
	if nok {
		Report("ERROR_ELEMENT_ID_REUSE", ref)
	}
	db.Elements[ref] = element
}

func (db *DbTopLevel) initElements() {
	for _, e := range db.Days {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Hours {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Teachers {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Subjects {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Rooms {
		db.testElement(e.Id, e)
	}
	for _, e := range db.RoomGroups {
		db.testElement(e.Id, e)
	}
	for _, e := range db.RoomChoiceGroups {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Groups {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Classes {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Courses {
		db.testElement(e.Id, e)
	}
	for _, e := range db.SuperCourses {
		db.testElement(e.Id, e)
	}
	for _, e := range db.SubCourses {
		db.testElement(e.Id, e)
	}
	for _, e := range db.Lessons {
		db.testElement(e.Id, e)
	}
	//TODO: Handle Constraints
}
