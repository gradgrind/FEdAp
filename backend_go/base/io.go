package base

import (
	"encoding/json"
	"io"
	"os"
)

func (db *DbTopLevel) SaveDb(fpath string) bool {
	// Save as JSON
	j, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		Report(`<Bug>Saving JSON: %v>`, err)
		return false
	}
	if err := os.WriteFile(fpath, j, 0666); err != nil {
		Report(`<Error>Saving file: %v>`, err)
		return false
	}
	return true
}

func LoadDb(fpath string) *DbTopLevel {
	// Open the  JSON file
	jsonFile, err := os.Open(fpath)
	if err != nil {
		Report(`<Error>Opening file: %v>`, err)
		return nil
	}
	// Remember to close the file at the end of the function
	defer jsonFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	Report(`<Info>Reading file: %s>`, fpath)
	v := NewDb()
	err = json.Unmarshal(byteValue, v)
	if err != nil {
		Report(`<Error>Invalid JSON: %s>`, err)
		return nil
	}
	v.initElements()
	return v
}

func (db *DbTopLevel) testElement(ref Ref, element Elem) {
	if ref == "" {
		Report(`<Error>Element has no Id:\n  -- %+v>`, element)
		db.SetInvalid()
	}
	_, nok := db.Elements[ref]
	if nok {
		Report(idduplicate, ref)
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
